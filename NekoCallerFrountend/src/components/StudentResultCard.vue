<script setup>
import { computed } from 'vue'
import { eventBadge } from '@/constants/events'

/** 被点名学生的基础信息 + 本次生效的事件标签。 */
const props = defineProps({
  student: { type: Object, required: true },
  eventType: { type: Number, default: 0 }
})

const badge = computed(() => eventBadge(props.eventType))
const enrollment = computed(() => props.student.enrollment_info || {})
</script>

<template>
  <div class="student-card">
    <el-descriptions :column="2" border>
      <el-descriptions-item label="学号">{{ student.student_info.student_id }}</el-descriptions-item>
      <el-descriptions-item label="姓名">
        <el-tag size="large" type="primary">{{ student.student_info.name }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="专业" :span="2">
        {{ student.student_info.major || '-' }}
      </el-descriptions-item>

      <el-descriptions-item label="当前积分">
        <el-statistic :value="enrollment.total_points" :precision="1">
          <template #suffix>分</template>
        </el-statistic>
      </el-descriptions-item>
      <el-descriptions-item label="点名次数">
        <el-statistic :value="enrollment.call_count">
          <template #suffix>次</template>
        </el-statistic>
      </el-descriptions-item>

      <el-descriptions-item label="转移权">
        <el-tag :type="enrollment.transfer_rights > 0 ? 'success' : 'info'">
          {{ enrollment.transfer_rights }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="跳过权">
        <el-tag :type="enrollment.skip_rights > 0 ? 'success' : 'info'">
          {{ enrollment.skip_rights }}
        </el-tag>
      </el-descriptions-item>

      <el-descriptions-item label="事件">
        <el-tooltip v-if="badge" effect="dark" placement="top" :content="badge.tip">
          <el-tag :type="badge.type" :color="badge.color" effect="dark">{{ badge.label }}</el-tag>
        </el-tooltip>
        <el-tag v-else type="info">无</el-tag>
      </el-descriptions-item>
    </el-descriptions>
  </div>
</template>
