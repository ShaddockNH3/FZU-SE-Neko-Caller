<script setup>
import { NButton, NCard, NForm, NFormItem, NIcon, NRadio, NRadioGroup, NSelect, NText } from 'naive-ui'
import { PhArrowCounterClockwise, PhPhoneCall } from '@phosphor-icons/vue'
import EventSelector from './EventSelector.vue'

/** 左侧配置面板：班级、模式、事件、开始与重置。状态全部在父级。 */
const props = defineProps({
  classes: { type: Array, default: () => [] },
  classId: { type: String, default: '' },
  mode: { type: Number, default: 0 },
  eventType: { type: Number, default: 0 },
  calling: { type: Boolean, default: false },
  resetting: { type: Boolean, default: false }
})

const emit = defineEmits([
  'update:classId',
  'update:mode',
  'update:eventType',
  'class-change',
  'mode-change',
  'start',
  'reset'
])

const MODE_OPTIONS = [
  { label: '随机', value: 0 },
  { label: '顺序', value: 1 },
  { label: '逆序', value: 2 },
  { label: '低分优先', value: 3 }
]

function classOptions() {
  return props.classes.map((item) => ({ label: item.class_name, value: item.class_id }))
}

function changeClass(value) {
  emit('update:classId', value)
  emit('class-change')
}

function changeMode(value) {
  emit('update:mode', value)
  emit('mode-change', value)
}
</script>

<template>
  <n-card size="small" title="点名配置">
    <n-form label-placement="left" :label-width="72" :show-feedback="false">
      <n-form-item label="班级">
        <n-select
          :value="props.classId"
          :options="classOptions()"
          placeholder="选择班级"
          @update:value="changeClass"
        />
      </n-form-item>

      <n-form-item label="模式">
        <n-radio-group :value="props.mode" @update:value="changeMode">
          <n-radio v-for="item in MODE_OPTIONS" :key="item.value" :value="item.value">
            {{ item.label }}
          </n-radio>
        </n-radio-group>
      </n-form-item>

      <n-form-item label="事件">
        <event-selector
          :model-value="props.eventType"
          :disabled="props.mode !== 0"
          @update:model-value="(value) => emit('update:eventType', value)"
        />
      </n-form-item>

      <n-text v-if="props.mode !== 0" depth="3" class="hint">
        事件只在随机模式下生效
      </n-text>

      <div class="actions">
        <n-button
          type="primary"
          size="large"
          block
          :disabled="!props.classId"
          :loading="props.calling"
          @click="emit('start')"
        >
          <template #icon>
            <n-icon><PhPhoneCall weight="fill" /></n-icon>
          </template>
          开始点名
        </n-button>

        <n-button
          v-if="props.mode === 1 || props.mode === 2"
          secondary
          size="large"
          block
          :disabled="!props.classId"
          :loading="props.resetting"
          @click="emit('reset')"
        >
          <template #icon>
            <n-icon><PhArrowCounterClockwise /></n-icon>
          </template>
          重置点名状态
        </n-button>

        <n-text v-if="props.mode === 1 || props.mode === 2" depth="3" class="hint">
          顺序 / 逆序点完一轮后，从这里重新开始
        </n-text>
      </div>
    </n-form>
  </n-card>
</template>

<style lang="scss" scoped>
.hint {
  display: block;
  font-size: 12px;
  margin: -4px 0 10px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 6px;
}
</style>
