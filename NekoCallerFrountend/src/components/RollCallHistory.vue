<script setup>
/** 点名历史：只展示，数据在 composable 里。 */
defineProps({
  records: { type: Array, default: () => [] }
})

function tagType(answerType) {
  if (answerType === 0) return 'success'
  if (answerType === 1) return 'warning'
  return 'info'
}
</script>

<template>
  <el-card class="history-card">
    <template #header>
      <span>点名历史</span>
    </template>

    <el-empty v-if="records.length === 0" description="还没有点名记录" :image-size="60" />

    <el-timeline v-else>
      <el-timeline-item
        v-for="(record, index) in records"
        :key="index"
        :timestamp="record.timestamp"
        placement="top"
      >
        <el-tag :type="tagType(record.answerType)" size="small">{{ record.studentName }}</el-tag>
        <span class="score">
          {{ record.scoreChange >= 0 ? '+' : '' }}{{ record.scoreChange }}
        </span>
      </el-timeline-item>
    </el-timeline>
  </el-card>
</template>

<style scoped>
.history-card {
  margin-top: 20px;
}

.score {
  margin-left: 10px;
}
</style>
