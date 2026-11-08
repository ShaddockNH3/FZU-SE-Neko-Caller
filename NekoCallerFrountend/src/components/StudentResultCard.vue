<script setup>
import { computed } from 'vue'
import { NDescriptions, NDescriptionsItem, NTag, NTooltip } from 'naive-ui'
import { eventBadge, naiveTagType } from '@/constants/events'

/** 被点名学生的基础信息 + 本次生效的事件标签。 */
const props = defineProps({
  student: { type: Object, required: true },
  eventType: { type: Number, default: 0 }
})

const badge = computed(() => eventBadge(props.eventType))
const enrollment = computed(() => props.student.enrollment_info || {})
const info = computed(() => props.student.student_info || {})
</script>

<template>
  <n-descriptions :column="2" bordered size="small" label-placement="left">
    <n-descriptions-item label="学号">{{ info.student_id }}</n-descriptions-item>

    <n-descriptions-item label="姓名">
      <n-tag type="primary" size="large" :bordered="false">{{ info.name }}</n-tag>
    </n-descriptions-item>

    <n-descriptions-item label="专业" :span="2">{{ info.major || '—' }}</n-descriptions-item>

    <n-descriptions-item label="当前积分">
      <span class="value">{{ Number(enrollment.total_points || 0).toFixed(1) }}</span>
    </n-descriptions-item>

    <n-descriptions-item label="点名次数">
      <span class="value">{{ enrollment.call_count || 0 }}</span>
    </n-descriptions-item>

    <n-descriptions-item label="转移权">
      <n-tag :type="enrollment.transfer_rights > 0 ? 'success' : 'default'" size="small" :bordered="false">
        {{ enrollment.transfer_rights || 0 }}
      </n-tag>
    </n-descriptions-item>

    <n-descriptions-item label="跳过权">
      <n-tag :type="enrollment.skip_rights > 0 ? 'success' : 'default'" size="small" :bordered="false">
        {{ enrollment.skip_rights || 0 }}
      </n-tag>
    </n-descriptions-item>

    <n-descriptions-item label="事件" :span="2">
      <n-tooltip v-if="badge" trigger="hover">
        <template #trigger>
          <n-tag :type="naiveTagType(badge.type)" size="small" :bordered="false">{{ badge.label }}</n-tag>
        </template>
        {{ badge.tip }}
      </n-tooltip>

      <n-tag v-else size="small" :bordered="false">无</n-tag>
    </n-descriptions-item>
  </n-descriptions>
</template>

<style lang="scss" scoped>
.value {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  font-variant-numeric: tabular-nums;
}
</style>
