<template>
  <div class="leaderboard-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>积分排行榜</span>
          <div class="filters">
            <el-select v-model="selectedClassId" placeholder="选择班级" style="width: 200px; margin-right: 10px;" @change="loadLeaderboard">
              <el-option v-for="cls in classList" :key="cls.class_id" :label="cls.class_name" :value="cls.class_id" />
            </el-select>
            <el-input-number v-model="topN" :min="1" :max="100" style="width: 120px; margin-right: 10px;" @change="loadLeaderboard" />
            <el-button type="primary" @click="loadLeaderboard">刷新</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20" style="margin-bottom: 20px;" v-if="stats">
        <el-col :span="8">
          <el-statistic title="总学生数" :value="stats.total_students">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="8">
          <el-statistic title="平均积分" :value="stats.average_points" :precision="2">
            <template #prefix>
              <el-icon><TrendCharts /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="8">
          <el-statistic title="累计点名" :value="stats.total_calls">
            <template #prefix>
              <el-icon><ChatDotSquare /></el-icon>
            </template>
          </el-statistic>
        </el-col>
      </el-row>

      <el-table :data="leaderboard" style="width: 100%" v-loading="loading">
        <el-table-column label="排名" width="80">
          <template #default="{ $index }">
            <el-tag v-if="$index === 0" type="danger" effect="dark">🥇 {{ $index + 1 }}</el-tag>
            <el-tag v-else-if="$index === 1" type="warning" effect="dark">🥈 {{ $index + 1 }}</el-tag>
            <el-tag v-else-if="$index === 2" type="success" effect="dark">🥉 {{ $index + 1 }}</el-tag>
            <span v-else style="font-weight: bold;">{{ $index + 1 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="student_id" label="学号" width="150" />
        <el-table-column prop="name" label="姓名" width="150" />
        <el-table-column label="专业">
          <template #default="{ row }">
            {{ row.major || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="total_points" label="总积分" width="120" sortable>
          <template #default="{ row }">
            <el-tag type="success" effect="plain">{{ row.total_points }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="call_count" label="点名次数" width="120" sortable />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { classAPI } from '@/api'
import { ElMessage } from 'element-plus'

const classList = ref([])
const selectedClassId = ref('')
const topN = ref(10)
const leaderboard = ref([])
const stats = ref(null)
const loading = ref(false)

const loadClasses = async () => {
  try {
    classList.value = await classAPI.getList()
    if (classList.value.length > 0) {
      selectedClassId.value = classList.value[0].class_id
      await loadLeaderboard()
    }
  } catch (error) {
    ElMessage.error('加载班级列表失败')
  }
}

const loadLeaderboard = async () => {
  if (!selectedClassId.value) return
  loading.value = true
  try {
    const [leaderboardData, statsData] = await Promise.all([
      classAPI.getLeaderboard(selectedClassId.value, topN.value),
      classAPI.getStats(selectedClassId.value)
    ])
    leaderboard.value = leaderboardData
    stats.value = statsData
  } catch (error) {
    ElMessage.error('加载排行榜失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadClasses()
})
</script>

<style scoped>
.leaderboard-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filters {
  display: flex;
  align-items: center;
}
</style>
