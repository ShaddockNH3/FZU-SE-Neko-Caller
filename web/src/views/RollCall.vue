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
              <el-radio-group v-model="mode">
                <el-radio :label="0">随机点名</el-radio>
                <el-radio :label="1">顺序点名</el-radio>
                <el-radio :label="2">逆序点名</el-radio>
                <el-radio :label="3">低分优先</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="事件类型">
              <el-radio-group v-model="eventType">
                <el-radio :label="0">无事件</el-radio>
                <el-radio :label="1">双倍积分</el-radio>
                <el-radio :label="2">疯狂星期四</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" size="large" style="width: 100%;" @click="startRollCall" :disabled="!selectedClassId" :loading="calling">
                <el-icon><PhoneFilled /></el-icon>
                开始点名
              </el-button>
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
              <el-descriptions-item label="事件">
                <el-tag v-if="eventType === 1" type="warning" effect="dark">双倍积分</el-tag>
                <el-tag v-else-if="eventType === 2" type="danger" effect="dark">疯狂星期四</el-tag>
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
                <el-radio :label="1">请求帮助</el-radio>
                <el-radio :label="2">跳过</el-radio>
                <el-radio :label="3">转移权</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="自定义分数" v-if="answerType === 0">
              <el-slider v-model="customScore" :min="0.5" :max="3" :step="0.5" show-stops :marks="{ 0.5: '0.5', 1: '1', 1.5: '1.5', 2: '2', 2.5: '2.5', 3: '3' }" />
              <div style="margin-top: 10px;">
                <el-tag>基础分：{{ customScore }}</el-tag>
                <el-tag v-if="eventType === 1" type="warning" style="margin-left: 10px;">双倍后：{{ customScore * 2 }}</el-tag>
                <el-tag v-if="eventType === 2" type="danger" style="margin-left: 10px;">疯四后：{{ (customScore * (Math.random() * 2 + 1)).toFixed(1) }}</el-tag>
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
import { ref, onMounted } from 'vue'
import { classAPI, rollCallAPI, rosterAPI } from '@/api'
import { ElMessage } from 'element-plus'

const classList = ref([])
const selectedClassId = ref('')
const mode = ref(0)
const eventType = ref(0)
const calling = ref(false)
const currentStudent = ref(null)
const answerType = ref(0)
const customScore = ref(1)
const targetStudentId = ref('')
const solving = ref(false)
const classRoster = ref([])
const history = ref([])

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
    // 后端返回的是 { base_response, roster_item }
    if (result.base_response?.code === 100 && result.roster_item) {
      currentStudent.value = result.roster_item
      answerType.value = 0
      customScore.value = 1
      targetStudentId.value = ''
      ElMessage.success(`点到了 ${result.roster_item.student_info.name}`)
    } else {
      throw new Error(result.base_response?.message || '点名失败')
    }
  } catch (error) {
    ElMessage.error(error.message || '点名失败')
  } finally {
    calling.value = false
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

  solving.value = true
  try {
    const payload = {
      enrollment_id: currentStudent.value.enrollment_info.enrollment_id,
      answer_type: answerType.value,
      event_type: eventType.value
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
    
    // 计算积分变化
    let scoreChange = 1.0
    if (answerType.value === 0) scoreChange = customScore.value
    else if (answerType.value === 1) scoreChange = 0.5
    else if (answerType.value === 2) scoreChange = -1.0
    else if (answerType.value === 3) scoreChange = -0.5
    
    // 应用事件加成
    if (answerType.value !== 3) {
      if (eventType.value === 1) scoreChange *= 2
      else if (eventType.value === 2) scoreChange *= 1.5
    }
    
    // 添加到历史记录
    history.value.unshift({
      timestamp: new Date().toLocaleTimeString(),
      studentName: currentStudent.value.student_info.name,
      answerType: answerType.value,
      scoreChange: scoreChange
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
