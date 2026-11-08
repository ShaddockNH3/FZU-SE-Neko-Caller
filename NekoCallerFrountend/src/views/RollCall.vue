<script setup>
import { NCard, NEmpty, NGi, NGrid, NIcon } from 'naive-ui'
import { PhHandPointing } from '@phosphor-icons/vue'
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
  <div>
    <header class="page-head">
      <h2>点名</h2>
      <p>选班级、挑模式，点一下开始；结算完自动进排行榜</p>
    </header>

    <n-grid cols="1 m:10" responsive="screen" :x-gap="16" :y-gap="16" item-responsive>
      <n-gi span="1 m:5">
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
      </n-gi>

      <n-gi span="1 m:5">
        <n-card v-if="!currentStudent" size="small" title="被点名学生">
          <n-empty description="点「开始点名」抽一个学生" class="empty">
            <template #icon>
              <n-icon><PhHandPointing /></n-icon>
            </template>
          </n-empty>
        </n-card>

        <n-card v-else size="small" title="被点名学生">
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
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<style lang="scss" scoped>
.empty {
  padding: 60px 0;
}
</style>
