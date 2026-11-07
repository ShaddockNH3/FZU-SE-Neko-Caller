<script setup>
import { PhoneFilled, RefreshLeft } from '@element-plus/icons-vue'
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
  <el-card>
    <template #header>
      <span>点名配置</span>
    </template>

    <el-form label-width="100px">
      <el-form-item label="选择班级">
        <el-select
          :model-value="props.classId"
          placeholder="请选择班级"
          class="full-width"
          @update:model-value="changeClass"
        >
          <el-option
            v-for="cls in props.classes"
            :key="cls.class_id"
            :label="cls.class_name"
            :value="cls.class_id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="点名模式">
        <el-radio-group :model-value="props.mode" @update:model-value="changeMode">
          <el-radio :label="0">随机点名</el-radio>
          <el-radio :label="1">顺序点名</el-radio>
          <el-radio :label="2">逆序点名</el-radio>
          <el-radio :label="3">低分优先</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="事件类型">
        <event-selector
          :model-value="props.eventType"
          :disabled="props.mode !== 0"
          @update:model-value="(value) => emit('update:eventType', value)"
        />
        <el-text v-if="props.mode !== 0" type="info" size="small" class="event-hint">
          仅随机点名模式可选事件
        </el-text>
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          size="large"
          class="full-width"
          :disabled="!props.classId"
          :loading="props.calling"
          @click="emit('start')"
        >
          <el-icon><PhoneFilled /></el-icon>
          开始点名
        </el-button>
      </el-form-item>

      <el-form-item v-if="props.mode === 1 || props.mode === 2">
        <el-button
          type="warning"
          size="large"
          class="full-width"
          :disabled="!props.classId"
          :loading="props.resetting"
          @click="emit('reset')"
        >
          <el-icon><RefreshLeft /></el-icon>
          重置点名状态
        </el-button>
        <el-text type="info" size="small" class="reset-hint">用于顺序/逆序点名重新开始</el-text>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<style scoped>
.full-width {
  width: 100%;
}

.event-hint {
  margin-left: 10px;
}

.reset-hint {
  margin-top: 5px;
}
</style>
