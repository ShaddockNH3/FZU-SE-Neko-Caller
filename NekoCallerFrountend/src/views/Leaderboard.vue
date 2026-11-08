<script setup>
import { computed, h, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  NButton,
  NCard,
  NDataTable,
  NGi,
  NGrid,
  NIcon,
  NInputNumber,
  NSelect,
  NStatistic,
  NTag
} from 'naive-ui'
import { PhArrowsClockwise, PhChatCircleText, PhChartLineUp, PhUsers } from '@phosphor-icons/vue'
import { classAPI } from '@/api'
import { palette } from '@/theme'
import { useBreakpoints } from '@/composables/useBreakpoints'

const route = useRoute()
const { isCompact, isNarrow } = useBreakpoints()

const classList = ref([])
const selectedClassId = ref('')
const topN = ref(10)
const leaderboard = ref([])
const stats = ref(null)
const loading = ref(false)

const classOptions = computed(() =>
  classList.value.map((item) => ({ label: item.class_name, value: item.class_id }))
)

const rankMedal = (index) => ['🥇', '🥈', '🥉'][index] || `${index + 1}`

const columns = computed(() => {
  const base = [
    {
      title: '排名',
      key: 'rank',
      width: 76,
      render: (row, index) =>
        index < 3
          ? h('span', { style: 'font-size:16px' }, rankMedal(index))
          : h('span', { style: 'color:#8b95a7;font-variant-numeric:tabular-nums' }, String(index + 1))
    },
    { title: '学号', key: 'student_id', width: isCompact.value ? 108 : 140 },
    { title: '姓名', key: 'name', minWidth: isNarrow.value ? 76 : 110 },
    {
      title: '专业',
      key: 'major',
      minWidth: 150,
      render: (row) => row.major || '—'
    },
    {
      title: '总积分',
      key: 'total_points',
      width: isCompact.value ? 82 : 110,
      sorter: (a, b) => a.total_points - b.total_points,
      render: (row) =>
        h(NTag, { size: 'small', type: 'success', bordered: false }, { default: () => row.total_points })
    },
    {
      title: '点名次数',
      key: 'call_count',
      width: isCompact.value ? 92 : 110,
      sorter: (a, b) => a.call_count - b.call_count
    }
  ]

  return isCompact.value ? base.filter((column) => column.key !== 'major') : base
})

async function loadClasses() {
  try {
    classList.value = await classAPI.getList()
  } catch (error) {
    console.error('加载班级列表失败:', error)
    return
  }

  // 从班级页点「排行榜」过来时带着 classId，优先用它
  const fromQuery = route.query.classId
  selectedClassId.value =
    (fromQuery && classList.value.some((item) => item.class_id === fromQuery) ? fromQuery : '') ||
    classList.value[0]?.class_id ||
    ''
}

async function loadLeaderboard() {
  if (!selectedClassId.value) return

  loading.value = true
  try {
    const [board, detail] = await Promise.all([
      classAPI.getLeaderboard(selectedClassId.value, topN.value),
      classAPI.getStats(selectedClassId.value)
    ])
    leaderboard.value = board || []
    stats.value = detail
  } catch (error) {
    console.error('加载排行榜失败:', error)
  } finally {
    loading.value = false
  }
}

watch(selectedClassId, loadLeaderboard)
watch(topN, loadLeaderboard)

onMounted(async () => {
  await loadClasses()
  await loadLeaderboard()
})
</script>

<template>
  <div>
    <header class="page-head">
      <h2>排行榜</h2>
      <p>按积分排序，答对加分、答错扣分都在这里体现</p>
    </header>

    <n-card size="small">
      <template #header>
        <div class="filters">
          <n-select
            v-model:value="selectedClassId"
            :options="classOptions"
            placeholder="选择班级"
            class="select"
          />
          <n-input-number v-model:value="topN" :min="1" :max="100" class="topn" />
          <n-button :loading="loading" secondary @click="loadLeaderboard">
            <template #icon>
              <n-icon><PhArrowsClockwise /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </template>

      <n-grid v-if="stats" cols="1 s:3" responsive="screen" :x-gap="16" :y-gap="12">
        <n-gi>
          <div class="metric">
            <n-icon :size="18" :color="palette.primary"><PhUsers /></n-icon>
            <n-statistic label="总学生数" :value="stats.total_students" />
          </div>
        </n-gi>
        <n-gi>
          <div class="metric">
            <n-icon :size="18" :color="palette.success"><PhChartLineUp /></n-icon>
            <n-statistic label="平均积分" :value="stats.average_points" :precision="1" />
          </div>
        </n-gi>
        <n-gi>
          <div class="metric">
            <n-icon :size="18" :color="palette.gold"><PhChatCircleText /></n-icon>
            <n-statistic label="累计点名" :value="stats.total_calls" />
          </div>
        </n-gi>
      </n-grid>

      <n-data-table
        class="table"
        :columns="columns"
        :data="leaderboard"
        :loading="loading"
        :bordered="false"
        :row-key="(row) => row.student_id"
        :scroll-x="isNarrow ? 440 : isCompact ? 460 : 700"
        size="medium"
      />
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.filters {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.select {
  width: 220px;
}

.topn {
  width: 110px;
}

.metric {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: #f7f9fc;
  border-radius: 10px;
}

.table {
  margin-top: 16px;
}

@media (max-width: 720px) {
  .select {
    width: 100%;
  }

  .topn {
    width: 90px;
  }
}
</style>
