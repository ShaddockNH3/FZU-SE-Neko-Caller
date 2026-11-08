<script setup>
import { NRadio, NRadioGroup, NTooltip } from 'naive-ui'
import { EVENT_OPTIONS } from '@/constants/events'

/** 事件单选：文案和说明都来自 constants/events.js，加事件只改数据。 */
defineProps({
  modelValue: { type: Number, default: 0 },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <n-radio-group
    :value="modelValue"
    :disabled="disabled"
    style="display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:6px 8px"
    @update:value="(value) => emit('update:modelValue', value)"
  >
    <template v-for="option in EVENT_OPTIONS" :key="option.value">
      <n-tooltip v-if="option.tip" trigger="hover" placement="top" :delay="300">
        <template #trigger>
          <n-radio :value="option.value">{{ option.label }}</n-radio>
        </template>

        <div class="tip">
          <strong>{{ option.tip.title }}</strong>
          <p>{{ option.tip.body }}</p>
          <span>{{ option.tip.keywords }}</span>
        </div>
      </n-tooltip>

      <n-radio v-else :value="option.value">{{ option.label }}</n-radio>
    </template>
  </n-radio-group>
</template>

<style lang="scss" scoped>
.tip {
  max-width: 280px;

  p {
    margin-top: 4px;
  }

  span {
    display: inline-block;
    margin-top: 4px;
    color: #9db8f0;
  }
}
</style>
