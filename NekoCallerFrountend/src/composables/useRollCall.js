import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { classAPI, rollCallAPI, rosterAPI } from '@/api'
import { ANSWER_TYPE, scoreDelta, scorePreview } from '@/constants/events'

/**
 * 点名页的状态与用例：选班、点名、结算、重置、历史。
 *
 * 原来这些散在 RollCall.vue 的 setup 里，和模板、事件标签搅在 500 行里；
 * 抽出来之后页面只负责布局，规则在 constants/events.js。
 */
export function useRollCall() {
  const classList = ref([])
  const selectedClassId = ref('')
  const mode = ref(0)
  const eventType = ref(0)
  const actualEventType = ref(0) // 后端校验后真正生效的事件
  const calling = ref(false)
  const resetting = ref(false)

  const currentStudent = ref(null)
  const classRoster = ref([])
  const history = ref([])

  // 结算表单
  const answerType = ref(ANSWER_TYPE.NORMAL)
  const customScore = ref(0)
  const targetStudentId = ref('')
  const solving = ref(false)

  const hasTransferRight = computed(
    () => (currentStudent.value?.enrollment_info?.transfer_rights || 0) > 0
  )
  const hasSkipRight = computed(
    () => (currentStudent.value?.enrollment_info?.skip_rights || 0) > 0
  )

  /** 转移目标：同班其他学生。 */
  const transferTargets = computed(() => {
    const current = currentStudent.value?.enrollment_info?.enrollment_id
    return classRoster.value.filter((item) => item.enrollment_info.enrollment_id !== current)
  })

  const preview = computed(() => scorePreview(actualEventType.value, customScore.value))
  const pendingDelta = computed(() =>
    scoreDelta({
      answerType: answerType.value,
      eventType: actualEventType.value,
      customScore: customScore.value
    })
  )

  async function loadClasses() {
    try {
      classList.value = await classAPI.getList()
    } catch (error) {
      console.error('加载班级失败:', error)
      ElMessage.error('加载班级列表失败')
    }
  }

  function onModeChange(value) {
    // 事件只在随机点名下可选，切换模式时清掉
    if (value !== 0) {
      eventType.value = 0
    }
  }

  async function onClassChange() {
    currentStudent.value = null
    try {
      classRoster.value = await rosterAPI.getClassRoster(selectedClassId.value)
    } catch (error) {
      console.error('加载花名册失败:', error)
      ElMessage.error('加载班级花名册失败')
    }
  }

  async function startRollCall() {
    calling.value = true
    try {
      const result = await rollCallAPI.call({
        class_id: selectedClassId.value,
        mode: mode.value,
        event_type: eventType.value
      })

      if (result.base_response?.code !== 100 || !result.roster_item) {
        throw new Error(result.base_response?.message || '点名失败')
      }

      currentStudent.value = result.roster_item
      actualEventType.value = result.actual_event_type ?? 0
      answerType.value = ANSWER_TYPE.NORMAL
      customScore.value = 0
      targetStudentId.value = ''

      const message = result.base_response.message || `点到了 ${result.roster_item.student_info.name}`
      const isSpecial = ['1024福报', '质数', 'Lucky 7'].some((word) => message.includes(word))
      if (isSpecial) {
        ElMessage({ message, type: 'success', duration: 5000, showClose: true })
      } else {
        ElMessage.success(message)
      }
    } catch (error) {
      console.error('点名失败:', error)
      ElMessage.error(error.message || '点名失败')
    } finally {
      calling.value = false
    }
  }

  async function resetRollCall() {
    try {
      await ElMessageBox.confirm(
        '重置后，该班级所有学生的点名次数将归零，是否继续？',
        '确认重置',
        { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
      )
    } catch {
      return // 取消
    }

    resetting.value = true
    try {
      const result = await rollCallAPI.reset(selectedClassId.value)
      if (result.code !== 100) {
        throw new Error(result.message || '重置失败')
      }
      ElMessage.success('重置成功')
      currentStudent.value = null
      await onClassChange()
    } catch (error) {
      console.error('重置失败:', error)
      ElMessage.error(error.message || '重置失败')
    } finally {
      resetting.value = false
    }
  }

  async function submitSolve() {
    const enrollment = currentStudent.value?.enrollment_info
    if (!enrollment) return

    if (answerType.value === ANSWER_TYPE.TRANSFER) {
      if (!targetStudentId.value) {
        ElMessage.warning('请选择转移目标学生')
        return
      }
      if (!hasTransferRight.value) {
        ElMessage.warning('转移权不足')
        return
      }
    }

    if (answerType.value === ANSWER_TYPE.SKIP && !hasSkipRight.value) {
      ElMessage.warning('跳过权不足')
      return
    }

    const payload = {
      enrollment_id: enrollment.enrollment_id,
      answer_type: answerType.value,
      event_type: actualEventType.value
    }

    if (answerType.value === ANSWER_TYPE.NORMAL) {
      payload.custom_score = customScore.value
    }

    if (answerType.value === ANSWER_TYPE.TRANSFER) {
      const target = classRoster.value.find(
        (item) => item.student_info.student_id === targetStudentId.value
      )
      if (!target) {
        ElMessage.error('未找到目标学生')
        return
      }
      payload.target_enrollment_id = target.enrollment_info.enrollment_id
      payload.custom_score = customScore.value // 转移也要带积分
    }

    solving.value = true
    try {
      await rollCallAPI.solve(payload)

      history.value.unshift({
        timestamp: new Date().toLocaleTimeString(),
        studentName: currentStudent.value.student_info.name,
        answerType: answerType.value,
        scoreChange: pendingDelta.value
      })

      ElMessage.success('结算成功')
      currentStudent.value = null
      await onClassChange()
    } catch (error) {
      console.error('结算失败:', error)
      ElMessage.error(error.response?.data?.message || '结算失败')
    } finally {
      solving.value = false
    }
  }

  onMounted(loadClasses)

  return {
    // 状态
    classList,
    selectedClassId,
    mode,
    eventType,
    actualEventType,
    calling,
    resetting,
    currentStudent,
    classRoster,
    history,
    answerType,
    customScore,
    targetStudentId,
    solving,
    // 派生
    hasTransferRight,
    hasSkipRight,
    transferTargets,
    preview,
    pendingDelta,
    // 动作
    onModeChange,
    onClassChange,
    startRollCall,
    resetRollCall,
    submitSolve
  }
}
