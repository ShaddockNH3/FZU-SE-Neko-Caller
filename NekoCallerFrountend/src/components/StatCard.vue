<script setup>
import { NCard, NIcon, NSkeleton, NStatistic } from 'naive-ui'

/**
 * 首页那张指标卡。四张卡只有图标、文案、数值和颜色不同，模板只留一份。
 */
defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], default: 0 },
  icon: { type: Object, required: true },
  color: { type: String, default: '#1e40af' },
  loading: { type: Boolean, default: false }
})
</script>

<template>
  <n-card size="small">
    <div class="row">
      <div class="badge" :style="{ background: `${color}14`, color }">
        <n-icon :size="22"><component :is="icon" weight="duotone" /></n-icon>
      </div>

      <n-statistic v-if="!loading" :label="label" :value="value" />
      <div v-else class="skeleton">
        <n-skeleton text width="48px" />
        <n-skeleton text width="72px" height="24px" />
      </div>
    </div>
  </n-card>
</template>

<style lang="scss" scoped>
.row {
  display: flex;
  align-items: center;
  gap: 14px;
}

.badge {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.skeleton {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
</style>
