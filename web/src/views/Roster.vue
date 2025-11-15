<template>
  <div class="roster-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <el-button icon="ArrowLeft" @click="$router.back()">返回</el-button>
            <span style="margin-left: 20px; font-size: 18px; font-weight: bold;">
              {{ className }} - 班级花名册
            </span>
          </div>
          <el-button type="primary" @click="exportRoster">
            <el-icon><Download /></el-icon>
            导出花名册
          </el-button>
        </div>
      </template>

      <el-table :data="roster" style="width: 100%" v-loading="loading">
        <el-table-column prop="student_info.student_id" label="学号" width="120" />
        <el-table-column prop="student_info.name" label="姓名" width="120" />
        <el-table-column label="专业">
          <template #default="{ row }">
            {{ row.student_info.major || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="enrollment_info.total_points" label="总积分" width="100" sortable />
        <el-table-column prop="enrollment_info.call_count" label="点名次数" width="100" sortable />
        <el-table-column prop="enrollment_info.transfer_rights" label="转移权" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.enrollment_info.transfer_rights > 0" type="success">
              {{ row.enrollment_info.transfer_rights }}
            </el-tag>
            <span v-else>0</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="removeStudent(row)">
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { rosterAPI, classAPI } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

const classId = route.params.classId
const className = ref('')
const roster = ref([])
const loading = ref(false)

const loadRoster = async () => {
  loading.value = true
  try {
    const [rosterData, classData] = await Promise.all([
      rosterAPI.getClassRoster(classId),
      classAPI.getDetail(classId)
    ])
    roster.value = rosterData
    className.value = classData.class_name
  } catch (error) {
    ElMessage.error('加载花名册失败')
  } finally {
    loading.value = false
  }
}

const exportRoster = async () => {
  try {
    const response = await classAPI.exportExcel(classId)
    const blob = new Blob([response.data], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${className.value}-花名册.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

const removeStudent = (row) => {
  ElMessageBox.confirm(
    `确定要将"${row.student_info.name}"从班级中移除吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await rosterAPI.removeStudent(row.enrollment_info.enrollment_id)
      ElMessage.success('移除成功')
      loadRoster()
    } catch (error) {
      ElMessage.error('移除失败')
    }
  })
}

onMounted(() => {
  loadRoster()
})
</script>

<style scoped>
.roster-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
