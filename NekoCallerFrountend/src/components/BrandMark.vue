<script setup>
import { NIcon } from 'naive-ui'
import { PhCat } from '@phosphor-icons/vue'

/** 侧栏和抽屉共用的品牌块。 */
defineProps({
  compact: { type: Boolean, default: false }
})
</script>

<template>
  <div class="brand" :class="{ 'is-compact': compact }">
    <n-icon :size="26" color="#1e40af"><PhCat weight="duotone" /></n-icon>

    <div v-if="!compact" class="brand-text">
      <h1>Neko Caller</h1>
      <span>课堂点名系统</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 64px;
  padding: 0 18px;
  border-bottom: 1px solid #eef1f7;

  &.is-compact {
    justify-content: center;
    padding: 0;
  }
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;

  h1 {
    font-size: 15px;
    font-weight: 600;
    color: #1f2937;
  }

  span {
    font-size: 11px;
    color: #8b95a7;
    letter-spacing: 0.6px;
  }
}
</style>
