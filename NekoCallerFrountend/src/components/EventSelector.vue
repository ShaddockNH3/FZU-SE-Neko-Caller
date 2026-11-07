<script setup>
import { EVENT_OPTIONS } from '@/constants/events'

/** 事件单选：文案和说明都来自 constants/events.js，加事件只改数据。 */
defineProps({
  modelValue: { type: Number, default: 0 },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <el-radio-group
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="(value) => emit('update:modelValue', value)"
  >
    <template v-for="option in EVENT_OPTIONS" :key="option.value">
      <el-tooltip v-if="option.tip" effect="dark" placement="top" :show-after="300">
        <template #content>
          <div class="event-tip">
            <div><strong>{{ option.tip.title }}</strong></div>
            <div class="event-tip-body">{{ option.tip.body }}</div>
            <div class="event-tip-keywords">{{ option.tip.keywords }}</div>
          </div>
        </template>
        <el-radio :label="option.value">{{ option.label }}</el-radio>
      </el-tooltip>
      <el-radio v-else :label="option.value">{{ option.label }}</el-radio>
    </template>
  </el-radio-group>
</template>

<style scoped>
.event-tip {
  max-width: 300px;
}

.event-tip-body {
  margin-top: 5px;
}

.event-tip-keywords {
  margin-top: 5px;
  color: #409eff;
}
</style>
