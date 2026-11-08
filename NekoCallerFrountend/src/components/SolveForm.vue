<script setup>
import {
  NButton,
  NDivider,
  NForm,
  NFormItem,
  NIcon,
  NRadio,
  NRadioGroup,
  NSelect,
  NSlider,
  NTag,
  NText
} from 'naive-ui'
import { PhCheckCircle } from '@phosphor-icons/vue'
import { naiveTagType } from '@/constants/events'

/**
 * 回答结算表单。积分规则在 constants/events.js，这里只管收集输入。
 */
const props = defineProps({
  student: { type: Object, required: true },
  answerType: { type: Number, default: 0 },
  customScore: { type: Number, default: 0 },
  targetStudentId: { type: String, default: '' },
  targets: { type: Array, default: () => [] },
  preview: { type: Object, default: null },
  pendingDelta: { type: Number, default: 0 },
  solving: { type: Boolean, default: false },
  canTransfer: { type: Boolean, default: false },
  canSkip: { type: Boolean, default: false }
})

const emit = defineEmits([
  'update:answerType',
  'update:customScore',
  'update:targetStudentId',
  'submit'
])

const SCORE_MARKS = {
  '-1': '-1',
  '-0.5': '-0.5',
  0: '0',
  0.5: '0.5',
  1: '1',
  1.5: '1.5',
  2: '2',
  2.5: '2.5',
  3: '3'
}

const targetOptions = () =>
  props.targets.map((item) => ({
    label: `${item.student_info.name}（${item.student_info.student_id}）`,
    value: item.student_info.student_id
  }))
</script>

<template>
  <div class="solve">
    <n-divider title-placement="left">回答结算</n-divider>

    <n-form label-placement="left" :label-width="88" :show-feedback="false">
      <n-form-item label="回答类型">
        <div>
          <n-radio-group
            :value="props.answerType"
            @update:value="(value) => emit('update:answerType', value)"
          >
            <n-radio :value="0">正常回答</n-radio>
            <n-radio :value="1">请求帮助</n-radio>
            <n-radio :value="2" :disabled="!props.canSkip">跳过</n-radio>
            <n-radio :value="3" :disabled="!props.canTransfer">转移权</n-radio>
          </n-radio-group>

          <n-text v-if="props.answerType === 2 && !props.canSkip" type="warning" class="hint">
            跳过权不足
          </n-text>
          <n-text v-if="props.answerType === 3 && !props.canTransfer" type="warning" class="hint">
            转移权不足
          </n-text>
        </div>
      </n-form-item>

      <n-form-item v-if="props.answerType === 0 || props.answerType === 3" label="回答分数">
        <div class="score-block">
          <n-slider
            :value="props.customScore"
            :min="-1"
            :max="3"
            :step="0.5"
            :marks="SCORE_MARKS"
            :format-tooltip="(value) => value.toFixed(1)"
            @update:value="(value) => emit('update:customScore', value)"
          />

          <div class="preview">
            <n-text v-if="props.answerType === 3" depth="3" class="hint">
              给转移目标选积分
            </n-text>

            <n-tag
              v-if="props.preview"
              size="small"
              :type="naiveTagType(props.preview.type)"
              :bordered="false"
            >
              {{ props.preview.text }}
            </n-tag>

            <n-text depth="3" class="hint">本次结算 {{ props.pendingDelta }} 分</n-text>
          </div>
        </div>
      </n-form-item>

      <n-form-item v-if="props.answerType === 3" label="转移目标">
        <div class="score-block">
          <n-select
            :value="props.targetStudentId"
            :options="targetOptions()"
            placeholder="选择转移目标"
            :disabled="!props.canTransfer"
            @update:value="(value) => emit('update:targetStudentId', value)"
          />
          <n-text v-if="!props.canTransfer" type="warning" class="hint">转移权不足</n-text>
        </div>
      </n-form-item>

      <n-button type="primary" size="large" block :loading="props.solving" @click="emit('submit')">
        <template #icon>
          <n-icon><PhCheckCircle weight="fill" /></n-icon>
        </template>
        提交结算
      </n-button>
    </n-form>
  </div>
</template>

<style lang="scss" scoped>
.solve {
  margin-top: 8px;
  padding-top: 4px;
}

.score-block {
  width: 100%;
  padding-top: 4px;
}

.hint {
  font-size: 12px;
}

.preview {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin: 14px 0 6px;
}
</style>
