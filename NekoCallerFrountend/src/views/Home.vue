<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NCard, NDataTable, NGi, NGrid, NIcon, NList, NListItem } from 'naive-ui'
import {
  PhChartLineUp,
  PhGraduationCap,
  PhMegaphone,
  PhTrophy,
  PhUsers
} from '@phosphor-icons/vue'
import StatCard from '@/components/StatCard.vue'
import { classAPI, studentAPI } from '@/api'

const router = useRouter()

const loading = ref(true)
const stats = ref({
  totalClasses: 0,
  totalStudents: 0,
  totalCalls: 0,
  averagePoints: 0
})
const recentClasses = ref([])

const statCards = computed(() => [
  { label: '班级总数', value: stats.value.totalClasses, icon: PhGraduationCap, color: '#1e40af' },
  { label: '学生总数', value: stats.value.totalStudents, icon: PhUsers, color: '#2f9e6f' },
  { label: '累计点名', value: stats.value.totalCalls, icon: PhMegaphone, color: '#c8963e' },
  { label: '平均积分', value: stats.value.averagePoints, icon: PhChartLineUp, color: '#6b5bd2' }
])

/**
 * 累计点名和人均积分都是算出来的，不是编的。
 *
 * 之前那版是「学生数 × 5 当累计点名、再乘 0.1 当今日」，后端摆着真数据不用。
 * 这里按班级拉一次排行榜（每项带 call_count 和 total_points），加起来就是真值。
 * 班级多的时候是 N 次请求，几十个班以内没问题，再多就该让后端加个汇总接口。
 */
async function loadData() {
  loading.value = true
  try {
    const [classes, students] = await Promise.all([classAPI.getList(), studentAPI.getList()])

    stats.value.totalClasses = classes.length
    stats.value.totalStudents = students.length
    recentClasses.value = classes.slice(0, 5)

    const boards = await Promise.all(
      classes.map((item) => classAPI.getLeaderboard(item.class_id, 999).catch(() => []))
    )

    let calls = 0
    let points = 0
    for (const board of boards) {
      for (const row of board || []) {
        calls += row.call_count || 0
        points += row.total_points || 0
      }
    }

    stats.value.totalCalls = calls
    stats.value.averagePoints = students.length ? Number((points / students.length).toFixed(1)) : 0
  } catch (error) {
    // 提示由拦截器统一弹
    console.error('加载首页数据失败:', error)
  } finally {
    loading.value = false
  }
}

function goRoster(classId) {
  router.push(`/roster/${classId}`)
}

function goRollCall(classId) {
  router.push({ path: '/roll-call', query: { classId } })
}

const columns = [
  { title: '班级名称', key: 'class_name', minWidth: 200 },
  {
    title: '学生人数',
    key: 'student_count',
    width: 120,
    render: (row) => row.student_ids?.length || 0
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row) => [
      h(
        NButton,
        { size: 'tiny', secondary: true, onClick: () => goRoster(row.class_id) },
        { default: () => '花名册' }
      ),
      h(
        NButton,
        { size: 'tiny', type: 'primary', style: 'margin-left:8px', onClick: () => goRollCall(row.class_id) },
        { default: () => '点名' }
      )
    ]
  }
]

onMounted(loadData)
</script>

<template>
  <div>
    <header class="page-head">
      <h2>首页</h2>
      <p>课堂点名的总览与快捷入口</p>
    </header>

    <n-grid cols="1 s:2 m:4" responsive="screen" :x-gap="16" :y-gap="16" item-responsive>
      <n-gi v-for="card in statCards" :key="card.label">
        <stat-card v-bind="card" :loading="loading" />
      </n-gi>
    </n-grid>

    <n-grid cols="1 m:2" responsive="screen" :x-gap="16" :y-gap="16" item-responsive class="row-gap">
      <n-gi>
        <n-card size="small" title="快速操作">
          <div class="quick">
            <n-button type="primary" secondary @click="router.push('/classes')">
              <template #icon>
                <n-icon><PhGraduationCap /></n-icon>
              </template>
              导入班级数据
            </n-button>

            <n-button type="primary" @click="router.push('/roll-call')">
              <template #icon>
                <n-icon><PhMegaphone /></n-icon>
              </template>
              开始点名
            </n-button>

            <n-button secondary @click="router.push('/leaderboard')">
              <template #icon>
                <n-icon><PhTrophy /></n-icon>
              </template>
              查看排行榜
            </n-button>
          </div>
        </n-card>
      </n-gi>

      <n-gi>
        <n-card size="small" title="这里能做什么">
          <n-list hoverable clickable :show-divider="false">
            <n-list-item v-for="line in [
              '随机 / 顺序 / 逆序 / 低分优先四种点名模式',
              '答对答错加减分，积分进排行榜',
              '随机事件：双倍积分、疯狂星期四、1024',
              '转移权和跳过权，答不上来能把机会递给别人'
            ]" :key="line">
              {{ line }}
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>

    <n-card size="small" title="最近班级" class="row-gap">
      <template #header-extra>
        <n-button text type="primary" @click="router.push('/classes')">查看全部</n-button>
      </template>

      <n-data-table
        :columns="columns"
        :data="recentClasses"
        :loading="loading"
        :bordered="false"
        :row-key="(row) => row.class_id"
        size="small"
      />
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.row-gap {
  margin-top: 16px;
}

.quick {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding: 6px 0;
}
</style>
