<template>
  <div class="roll-call-container">
    <el-row :gutter="20">
      <!-- 左侧配置区 -->
      <el-col :span="10">
        <el-card>
          <template #header>
            <span>点名配置</span>
          </template>

          <el-form label-width="100px">
            <el-form-item label="选择班级">
              <el-select v-model="selectedClassId" placeholder="请选择班级" style="width: 100%;" @change="onClassChange">
                <el-option v-for="cls in classList" :key="cls.class_id" :label="cls.class_name" :value="cls.class_id" />
              </el-select>
            </el-form-item>

            <el-form-item label="点名模式">
              <el-radio-group v-model="mode" @change="onModeChange">
                <el-radio :label="0">随机点名</el-radio>
                <el-radio :label="1">顺序点名</el-radio>
                <el-radio :label="2">逆序点名</el-radio>
                <el-radio :label="3">低分优先</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="事件类型">
              <el-radio-group v-model="eventType" :disabled="mode !== 0">
                <el-radio :label="0">无事件</el-radio>
                
                <el-tooltip effect="dark" placement="top" :show-after="300">
                  <template #content>
                    <div style="max-width: 300px;">
                      <div><strong>双倍积分 (Double Point)</strong></div>
                      <div style="margin-top: 5px;">触发简单的幸运事件，被点名同学若回答正确，本次所获基础积分翻倍。</div>
                      <div style="margin-top: 5px; color: #409eff;">Keywords: 运气，翻倍奖励</div>
                    </div>
                  </template>
                  <el-radio :label="1">双倍积分</el-radio>
                </el-tooltip>
                
                <el-tooltip effect="dark" placement="top" :show-after="300">
                  <template #content>
                    <div style="max-width: 300px;">
                      <div><strong>疯狂星期四 (Crazy Thursday)</strong></div>
                      <div style="margin-top: 5px;">日期限定彩蛋。检测系统时间是否为星期四，触发"V我50"相关机制（如积分为50因数的同学权重增加，或获得特殊Buff）。</div>
                      <div style="margin-top: 5px; color: #409eff;">Keywords: 日期检测，玩梗，特定人群Buff</div>
                    </div>
                  </template>
                  <el-radio :label="2">疯狂星期四</el-radio>
                </el-tooltip>
                
                <el-tooltip effect="dark" placement="top" :show-after="300">
                  <template #content>
                    <div style="max-width: 300px;">
                      <div><strong>1024 程序员福报 (Blessing 1024)</strong></div>
                      <div style="margin-top: 5px;">当点名发生的系统时间（秒）为 2 的幂次方（1, 2, 4, 8...）时触发。致敬二进制文化，回答积分固定为 1.024 分，且具备"免死金牌"效果（答错不扣分）。</div>
                      <div style="margin-top: 5px; color: #409eff;">Keywords: 二进制，免死金牌，极客文化</div>
                    </div>
                  </template>
                  <el-radio :label="3">1024福报</el-radio>
                </el-tooltip>
                
                <el-tooltip effect="dark" placement="top" :show-after="300">
                  <template #content>
                    <div style="max-width: 300px;">
                      <div><strong>质数的孤独 (Solitude of Primes)</strong></div>
                      <div style="margin-top: 5px;">当点名发生的系统时间（分）为质数（2, 3, 5, 7...）时触发。寓意"在软工课上你并不孤独"，给予额外的暖心积分加成（如 +0.37分）。</div>
                      <div style="margin-top: 5px; color: #409eff;">Keywords: 数学规律，额外津贴，人文关怀</div>
                    </div>
                  </template>
                  <el-radio :label="4">质数的孤独</el-radio>
                </el-tooltip>
                
                <el-tooltip effect="dark" placement="top" :show-after="300">
                  <template #content>
                    <div style="max-width: 300px;">
                      <div><strong>幸运 7 大奖 (Lucky 7 Jackpot)</strong></div>
                      <div style="margin-top: 5px;">基于学号的被动触发事件。若被选中的同学学号末尾字符为 '7'，则自动触发"大奖"特效。回答错误不扣分，回答正确获得幸运加分。</div>
                      <div style="margin-top: 5px; color: #409eff;">Keywords: 字符串匹配，豁免权，强运</div>
                    </div>
                  </template>
                  <el-radio :label="5">幸运7大奖</el-radio>
                </el-tooltip>
              </el-radio-group>
              <el-text v-if="mode !== 0" type="info" size="small" style="margin-left: 10px;">仅随机点名模式可选事件</el-text>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" size="large" style="width: 100%;" @click="startRollCall" :disabled="!selectedClassId" :loading="calling">
                <el-icon><PhoneFilled /></el-icon>
                开始点名
              </el-button>
            </el-form-item>

            <el-form-item v-if="mode === 1 || mode === 2">
              <el-button type="warning" size="large" style="width: 100%;" @click="resetRollCall" :disabled="!selectedClassId" :loading="resetting">
                <el-icon><RefreshLeft /></el-icon>
                重置点名状态
              </el-button>
              <el-text type="info" size="small" style="margin-top: 5px;">用于顺序/逆序点名重新开始</el-text>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 点名历史 -->
        <el-card style="margin-top: 20px;">
          <template #header>
            <span>点名历史</span>
          </template>
          <el-timeline>
            <el-timeline-item v-for="(record, index) in history" :key="index" :timestamp="record.timestamp" placement="top">
              <el-tag :type="record.answerType === 0 ? 'success' : record.answerType === 1 ? 'warning' : 'info'" size="small">
                {{ record.studentName }}
              </el-tag>
              <span style="margin-left: 10px;">{{ record.scoreChange >= 0 ? '+' : '' }}{{ record.scoreChange }}</span>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>

      <!-- 右侧结果区 -->
      <el-col :span="14">
        <el-card v-if="!currentStudent">
          <el-empty description="点击开始点名按钮开始" />
        </el-card>

        <el-card v-else>
          <template #header>
            <span>被点名学生</span>
          </template>

          <div class="student-card">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="学号">{{ currentStudent.student_info.student_id }}</el-descriptions-item>
              <el-descriptions-item label="姓名">
                <el-tag size="large" type="primary">{{ currentStudent.student_info.name }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="专业" :span="2">{{ currentStudent.student_info.major || '-' }}</el-descriptions-item>
              <el-descriptions-item label="当前积分">
                <el-statistic :value="currentStudent.enrollment_info.total_points" :precision="1">
                  <template #suffix>分</template>
                </el-statistic>
              </el-descriptions-item>
              <el-descriptions-item label="点名次数">
                <el-statistic :value="currentStudent.enrollment_info.call_count">
                  <template #suffix>次</template>
                </el-statistic>
              </el-descriptions-item>
              <el-descriptions-item label="转移权">
                <el-tag :type="currentStudent.enrollment_info.transfer_rights > 0 ? 'success' : 'info'">
                  {{ currentStudent.enrollment_info.transfer_rights }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="跳过权">
                <el-tag :type="currentStudent.enrollment_info.skip_rights > 0 ? 'success' : 'info'">
                  {{ currentStudent.enrollment_info.skip_rights }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="事件">
                <el-tooltip v-if="actualEventType === 1" effect="dark" placement="top" content="双倍积分: 回答正确积分翻倍">
                  <el-tag type="warning" effect="dark">双倍积分</el-tag>
                </el-tooltip>
                <el-tooltip v-else-if="actualEventType === 2" effect="dark" placement="top" content="疯狂星期四: V我50特殊机制">
                  <el-tag type="danger" effect="dark">疯狂星期四</el-tag>
                </el-tooltip>
                <el-tooltip v-else-if="actualEventType === 3" effect="dark" placement="top" content="1024福报: +1.024分且答错免扣分">
                  <el-tag type="success" effect="dark">1024福报</el-tag>
                </el-tooltip>
                <el-tooltip v-else-if="actualEventType === 4" effect="dark" placement="top" content="质数的孤独: 回答正确额外+0.37分">
                  <el-tag type="primary" effect="dark">质数的孤独</el-tag>
                </el-tooltip>
                <el-tooltip v-else-if="actualEventType === 5" effect="dark" placement="top" content="幸运7大奖: 答错不扣分">
                  <el-tag color="#f56c6c" effect="dark">幸运7大奖</el-tag>
                </el-tooltip>
                <el-tag v-else type="info">无</el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <!-- 结算表单 -->
          <el-divider>回答结算</el-divider>
          <el-form label-width="120px">
            <el-form-item label="回答类型">
              <el-radio-group v-model="answerType">
                <el-radio :label="0">正常回答</el-radio>
                <el-radio :label="1">请求帮助（准确重复）</el-radio>
                <el-radio :label="2" :disabled="!hasSkipRight">跳过（需跳过权）</el-radio>
                <el-radio :label="3" :disabled="!hasTransferRight">转移权</el-radio>
              </el-radio-group>
              <div style="margin-top: 5px;">
                <el-text v-if="answerType === 2 && !hasSkipRight" type="warning" size="small">跳过权不足</el-text>
                <el-text v-if="answerType === 3 && !hasTransferRight" type="warning" size="small">转移权不足</el-text>
              </div>
            </el-form-item>

            <el-form-item label="回答分数" v-if="answerType === 0">
              <el-slider v-model="customScore" :min="-1" :max="3" :step="0.5" show-stops :marks="{ '-1': '-1', '-0.5': '-0.5', 0: '0', 0.5: '0.5', 1: '1', 1.5: '1.5', 2: '2', 2.5: '2.5', 3: '3' }" />
              <div style="margin-top: 10px;">
                <el-tag v-if="actualEventType === 1" type="warning">双倍后：{{ (customScore + 1) * 2 }}</el-tag>
                <el-tag v-if="actualEventType === 2" type="danger" style="margin-left: 10px;">疯四后：{{ ((customScore + 1) * 1.5).toFixed(1) }}</el-tag>
                <el-tag v-if="actualEventType === 3" type="success" style="margin-left: 10px;">
                  1024福报：{{ customScore > 0 ? '1.024' : customScore < 0 ? '1.0(免扣分)' : (customScore + 1).toFixed(1) }}
                </el-tag>
                <el-tag v-if="actualEventType === 4" type="primary" style="margin-left: 10px;">
                  质数孤独：{{ customScore >= 0 ? ((customScore + 1) + 0.37).toFixed(2) : (customScore + 1).toFixed(1) }}
                </el-tag>
                <el-tag v-if="actualEventType === 5" color="#f56c6c" style="margin-left: 10px;">
                  幸运7：{{ customScore < 0 ? '1.0(免扣分)' : (customScore + 1).toFixed(1) }}
                </el-tag>
              </div>
            </el-form-item>

            <el-form-item label="转移目标" v-if="answerType === 3">
              <el-select v-model="targetStudentId" placeholder="选择转移目标学生" style="width: 100%;" :disabled="currentStudent.enrollment_info.transfer_rights <= 0">
                <el-option v-for="student in classRoster.filter(s => s.enrollment_info.enrollment_id !== currentStudent.enrollment_info.enrollment_id)" :key="student.enrollment_info.enrollment_id" :label="`${student.student_info.name} (${student.student_info.student_id})`" :value="student.student_info.student_id" />
              </el-select>
              <el-text v-if="currentStudent.enrollment_info.transfer_rights <= 0" type="warning" size="small" style="margin-top: 5px;">转移权不足</el-text>
            </el-form-item>

            <el-form-item>
              <el-button type="success" size="large" style="width: 100%;" @click="submitSolve" :loading="solving">
                <el-icon><CircleCheckFilled /></el-icon>
                提交结算
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { classAPI, rollCallAPI, rosterAPI } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const classList = ref([])
const selectedClassId = ref('')
const mode = ref(0)
const eventType = ref(0)
const actualEventType = ref(0) // 实际触发的事件类型（由后端验证后返回）
const calling = ref(false)
const resetting = ref(false)
const currentStudent = ref(null)
const answerType = ref(0)
const customScore = ref(0)
const targetStudentId = ref('')
const solving = ref(false)
const classRoster = ref([])
const history = ref([])

// 计算属性：是否有转移权
const hasTransferRight = computed(() => {
  return currentStudent.value?.enrollment_info?.transfer_rights > 0
})

// 计算属性：是否有跳过权
const hasSkipRight = computed(() => {
  return currentStudent.value?.enrollment_info?.skip_rights > 0
})

// 模式切换时重置事件类型
const onModeChange = () => {
  if (mode.value !== 0) {
    eventType.value = 0
  }
}

const loadClasses = async () => {
  try {
    classList.value = await classAPI.getList()
  } catch (error) {
    ElMessage.error('加载班级列表失败')
  }
}

const onClassChange = async () => {
  currentStudent.value = null
  try {
    classRoster.value = await rosterAPI.getClassRoster(selectedClassId.value)
  } catch (error) {
    ElMessage.error('加载班级花名册失败')
  }
}

const startRollCall = async () => {
  calling.value = true
  try {
    const result = await rollCallAPI.call({
      class_id: selectedClassId.value,
      mode: mode.value,
      event_type: eventType.value
    })
    // 后端返回的是 { base_response, roster_item, actual_event_type }
    if (result.base_response?.code === 100 && result.roster_item) {
      currentStudent.value = result.roster_item
      actualEventType.value = result.actual_event_type ?? 0 // 保存实际触发的事件类型
      answerType.value = 0
      customScore.value = 0
      targetStudentId.value = ''
      
      // 显示点名成功消息，包含可能的事件提示
      const message = result.base_response.message || `点到了 ${result.roster_item.student_info.name}`
      
      // 根据消息内容确定消息类型
      if (message.includes('1024福报') || message.includes('质数') || message.includes('Lucky 7')) {
        ElMessage({
          message: message,
          type: 'success',
          duration: 5000,
          showClose: true
        })
      } else {
        ElMessage.success(message)
      }
    } else {
      throw new Error(result.base_response?.message || '点名失败')
    }
  } catch (error) {
    ElMessage.error(error.message || '点名失败')
  } finally {
    calling.value = false
  }
}

const resetRollCall = async () => {
  try {
    await ElMessageBox.confirm(
      '重置后，该班级所有学生的点名次数将归零，是否继续？',
      '确认重置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    resetting.value = true
    const result = await rollCallAPI.reset(selectedClassId.value)
    
    if (result.code === 100) {
      ElMessage.success('重置成功')
      // 刷新花名册
      await onClassChange()
      // 清除当前点名学生
      currentStudent.value = null
    } else {
      throw new Error(result.message || '重置失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '重置失败')
    }
  } finally {
    resetting.value = false
  }
}

const submitSolve = async () => {
  if (answerType.value === 3 && !targetStudentId.value) {
    ElMessage.warning('请选择转移目标学生')
    return
  }

  // 检查转移权
  if (answerType.value === 3 && currentStudent.value.enrollment_info.transfer_rights <= 0) {
    ElMessage.warning('转移权不足')
    return
  }

  // 检查跳过权
  if (answerType.value === 2 && currentStudent.value.enrollment_info.skip_rights <= 0) {
    ElMessage.warning('跳过权不足')
    return
  }

  solving.value = true
  try {
    const payload = {
      enrollment_id: currentStudent.value.enrollment_info.enrollment_id,
      answer_type: answerType.value,
      event_type: actualEventType.value // 使用实际触发的事件类型
    }

    if (answerType.value === 0) {
      payload.custom_score = customScore.value
    }

    if (answerType.value === 3) {
      // 根据student_id找到对应的enrollment_id
      const targetStudent = classRoster.value.find(s => s.student_info.student_id === targetStudentId.value)
      if (!targetStudent) {
        ElMessage.error('未找到目标学生')
        return
      }
      payload.target_enrollment_id = targetStudent.enrollment_info.enrollment_id
    }

    const result = await rollCallAPI.solve(payload)
    
    // 使用后端返回的实际事件类型来计算积分
    const actualEvent = result.data?.actual_event_type ?? 0
    
    // 计算积分变化（根据新规则）
    let scoreChange = 0
    if (answerType.value === 0) {
      // 正常回答：到达+1 + 回答得分（-1到3）
      scoreChange = 1.0 + customScore.value
      
      // 应用特殊事件规则（使用实际触发的事件类型）
      if (actualEvent === 3) {
        // 1024福报
        if (customScore.value > 0) {
          scoreChange = 1.024 // 回答正确固定为1.024
        } else if (customScore.value < 0) {
          scoreChange = 1.0 // 回答错误免扣分
        }
      } else if (actualEvent === 4) {
        // 质数的孤独：回答正确额外+0.37
        if (customScore.value >= 0) {
          scoreChange += 0.37
        }
      } else if (actualEvent === 5) {
        // 幸运7：回答错误免扣分
        if (customScore.value < 0) {
          scoreChange = 1.0
        }
      } else if (actualEvent === 1) {
        // 双倍积分
        scoreChange *= 2
      } else if (actualEvent === 2) {
        // 疯狂星期四
        scoreChange *= 1.5
      }
    } else if (answerType.value === 1) {
      // 请求帮助：到达+1 + 准确重复+0.5
      scoreChange = 1.5
      
      // 质数的孤独对帮助也生效
      if (eventType.value === 4) {
        scoreChange += 0.37
      } else if (eventType.value === 1) {
        scoreChange *= 2
      } else if (eventType.value === 2) {
        scoreChange *= 1.5
      }
    } else if (answerType.value === 2) {
      // 跳过：到达+1
      scoreChange = 1.0
    } else if (answerType.value === 3) {
      // 转移：到达+1
      scoreChange = 1.0
    }
    
    // 添加到历史记录
    history.value.unshift({
      timestamp: new Date().toLocaleTimeString(),
      studentName: currentStudent.value.student_info.name,
      answerType: answerType.value,
      scoreChange: Math.round(scoreChange * 1000) / 1000
    })

    ElMessage.success('结算成功')
    currentStudent.value = null
    
    // 刷新花名册
    await onClassChange()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '结算失败')
  } finally {
    solving.value = false
  }
}

onMounted(() => {
  loadClasses()
})
</script>

<style scoped>
.roll-call-container {
  padding: 20px;
}

.student-card {
  margin-bottom: 20px;
}
</style>
