<script setup>
import { CircleCheckFilled } from '@element-plus/icons-vue'

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
</script>

<template>
  <div class="solve-panel">
    <el-divider>回答结算</el-divider>

    <el-form label-width="120px">
      <el-form-item label="回答类型">
        <el-radio-group
          :model-value="props.answerType"
          @update:model-value="(value) => emit('update:answerType', value)"
        >
          <el-radio :label="0">正常回答</el-radio>
          <el-radio :label="1">请求帮助（准确重复）</el-radio>
          <el-radio :label="2" :disabled="!props.canSkip">跳过（需跳过权）</el-radio>
          <el-radio :label="3" :disabled="!props.canTransfer">转移权</el-radio>
        </el-radio-group>

        <div class="hint">
          <el-text v-if="props.answerType === 2 && !props.canSkip" type="warning" size="small">
            跳过权不足
          </el-text>
          <el-text v-if="props.answerType === 3 && !props.canTransfer" type="warning" size="small">
            转移权不足
          </el-text>
        </div>
      </el-form-item>

      <el-form-item v-if="props.answerType === 0 || props.answerType === 3" label="回答分数">
        <el-slider
          :model-value="props.customScore"
          :min="-1"
          :max="3"
          :step="0.5"
          show-stops
          :marks="SCORE_MARKS"
          @update:model-value="(value) => emit('update:customScore', value)"
        />

        <div class="preview">
          <el-text v-if="props.answerType === 3" type="primary" size="small">
            为转移目标选择积分
          </el-text>
          <el-tag v-if="props.preview" :type="props.preview.type" :color="props.preview.color">
            {{ props.preview.text }}
          </el-tag>
          <el-text type="info" size="small">本次结算 {{ props.pendingDelta }} 分</el-text>
        </div>
      </el-form-item>

      <el-form-item v-if="props.answerType === 3" label="转移目标">
        <el-select
          :model-value="props.targetStudentId"
          placeholder="选择转移目标学生"
          class="full-width"
          :disabled="!props.canTransfer"
          @update:model-value="(value) => emit('update:targetStudentId', value)"
        >
          <el-option
            v-for="item in props.targets"
            :key="item.enrollment_info.enrollment_id"
            :label="`${item.student_info.name} (${item.student_info.student_id})`"
            :value="item.student_info.student_id"
          />
        </el-select>
        <el-text v-if="!props.canTransfer" type="warning" size="small" class="hint">
          转移权不足
        </el-text>
      </el-form-item>

      <el-form-item>
        <el-button
          type="success"
          size="large"
          class="full-width"
          :loading="props.solving"
          @click="emit('submit')"
        >
          <el-icon><CircleCheckFilled /></el-icon>
          提交结算
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped>
.full-width {
  width: 100%;
}

.hint {
  margin-top: 5px;
}

.preview {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
