<template>
  <div class="home-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#409EFF"><School /></el-icon>
            <div class="stat-info">
              <div class="stat-label">班级总数</div>
              <div class="stat-value">{{ stats.totalClasses }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#67C23A"><User /></el-icon>
            <div class="stat-info">
              <div class="stat-label">学生总数</div>
              <div class="stat-value">{{ stats.totalStudents }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#E6A23C"><Promotion /></el-icon>
            <div class="stat-info">
              <div class="stat-label">今日点名</div>
              <div class="stat-value">{{ stats.todayCalls }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#F56C6C"><TrophyBase /></el-icon>
            <div class="stat-info">
              <div class="stat-label">累计点名</div>
              <div class="stat-value">{{ stats.totalCalls }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/classes')">
              <el-icon><Upload /></el-icon>
              导入班级数据
            </el-button>
            <el-button type="success" size="large" @click="$router.push('/roll-call')">
              <el-icon><Promotion /></el-icon>
              开始点名
            </el-button>
            <el-button type="warning" size="large" @click="$router.push('/leaderboard')">
              <el-icon><TrophyBase /></el-icon>
              查看排行榜
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统说明</span>
            </div>
          </template>
          <div class="system-info">
            <h4>功能特性</h4>
            <ul>
              <li>✅ 支持Excel导入学生名单</li>
              <li>✅ 多种点名模式（随机、顺序、低分优先）</li>
              <li>✅ 积分管理与排行榜</li>
              <li>✅ 点名转移权机制</li>
              <li>✅ 随机事件系统（双倍积分、疯狂星期四）</li>
              <li>✅ 数据可视化统计</li>
            </ul>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近班级</span>
              <el-button text @click="$router.push('/classes')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentClasses" style="width: 100%">
            <el-table-column prop="class_name" label="班级名称" />
            <el-table-column label="学生人数">
              <template #default="{ row }">
                {{ row.student_ids?.length || 0 }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="$router.push(`/roster/${row.class_id}`)">
                  查看花名册
                </el-button>
                <el-button size="small" type="primary" @click="goRollCall(row.class_id)">
                  点名
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { classAPI, studentAPI } from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()

const stats = ref({
  totalClasses: 0,
  totalStudents: 0,
  todayCalls: 0,
  totalCalls: 0
})

const recentClasses = ref([])

const loadData = async () => {
  try {
    const [classes, students] = await Promise.all([
      classAPI.getList(),
      studentAPI.getList()
    ])
    
    stats.value.totalClasses = classes.length
    stats.value.totalStudents = students.length
    recentClasses.value = classes.slice(0, 5)
    
    // 计算累计点名次数（模拟数据）
    stats.value.totalCalls = classes.reduce((sum, cls) => {
      return sum + (cls.student_ids?.length || 0) * 5
    }, 0)
    stats.value.todayCalls = Math.floor(stats.value.totalCalls * 0.1)
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const goRollCall = (classId) => {
  router.push({ path: '/roll-call', query: { classId } })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.home-container {
  padding: 20px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 10px;
}

.stat-icon {
  margin-right: 20px;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-actions {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
}

.system-info h4 {
  margin-bottom: 15px;
  color: #303133;
}

.system-info ul {
  list-style: none;
  padding: 0;
}

.system-info li {
  padding: 8px 0;
  color: #606266;
  font-size: 14px;
}
</style>
