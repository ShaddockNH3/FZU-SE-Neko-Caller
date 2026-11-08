<script setup>
import { NCard, NEmpty, NTag, NTimeline, NTimelineItem } from 'naive-ui'

/** 点名历史：只展示，数据在 composable 里。 */
defineProps({
  records: { type: Array, default: () => [] }
})

function tagType(answerType) {
  if (answerType === 0) return 'success'
  if (answerType === 1) return 'warning'
  return 'default'
}
</script>

<template>
  <n-card size="small" title="点名历史" class="history">
    <n-empty v-if="records.length === 0" description="还没有点名记录" size="small" class="empty" />

    <n-timeline v-else>
      <n-timeline-item
        v-for="(record, index) in records"
        :key="index"
        :time="record.timestamp"
        type="success"
      >
        <n-tag :type="tagType(record.answerType)" size="small" :bordered="false">
          {{ record.studentName }}
        </n-tag>
        <span class="score" :class="{ minus: record.scoreChange < 0 }">
          {{ record.scoreChange >= 0 ? '+' : '' }}{{ record.scoreChange }}
        </span>
      </n-timeline-item>
    </n-timeline>
  </n-card>
</template>

<style lang="scss" scoped>
.history {
  margin-top: 16px;
}

.empty {
  padding: 20px 0;
}

.score {
  margin-left: 10px;
  font-variant-numeric: tabular-nums;
  color: #2f9e6f;

  &.minus {
    color: #c2453d;
  }
}
</style>
