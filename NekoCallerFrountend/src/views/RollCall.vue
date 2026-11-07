<script setup>
import { useRollCall } from '@/composables/useRollCall'
import RollCallConfig from '@/components/RollCallConfig.vue'
import RollCallHistory from '@/components/RollCallHistory.vue'
import StudentResultCard from '@/components/StudentResultCard.vue'
import SolveForm from '@/components/SolveForm.vue'

// 状态和用例都在 composable 里，这里只摆版面、连事件。
const {
  classList,
  selectedClassId,
  mode,
  eventType,
  actualEventType,
  calling,
  resetting,
  currentStudent,
  history,
  answerType,
  customScore,
  targetStudentId,
  solving,
  hasTransferRight,
  hasSkipRight,
  transferTargets,
  preview,
  pendingDelta,
  onModeChange,
  onClassChange,
  startRollCall,
  resetRollCall,
  submitSolve
} = useRollCall()
</script>

<template>
  <div class="roll-call-container">
    <el-row :gutter="20">
      <el-col :span="10">
        <roll-call-config
          v-model:class-id="selectedClassId"
          v-model:mode="mode"
          v-model:event-type="eventType"
          :classes="classList"
          :calling="calling"
          :resetting="resetting"
          @class-change="onClassChange"
          @mode-change="onModeChange"
          @start="startRollCall"
          @reset="resetRollCall"
        />

        <roll-call-history :records="history" />
      </el-col>

      <el-col :span="14">
        <el-card v-if="!currentStudent">
          <el-empty description="点击开始点名按钮开始" />
        </el-card>

        <el-card v-else>
          <template #header>
            <span>被点名学生</span>
          </template>

          <student-result-card :student="currentStudent" :event-type="actualEventType" />

          <solve-form
            v-model:answer-type="answerType"
            v-model:custom-score="customScore"
            v-model:target-student-id="targetStudentId"
            :student="currentStudent"
            :targets="transferTargets"
            :preview="preview"
            :pending-delta="pendingDelta"
            :solving="solving"
            :can-transfer="hasTransferRight"
            :can-skip="hasSkipRight"
            @submit="submitSolve"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.roll-call-container {
  padding: 4px;
}
</style>
