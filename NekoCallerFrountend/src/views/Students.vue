<template>
  <div class="students-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>学生管理</span>
          <el-input
            v-model="searchText"
            placeholder="搜索学生姓名或学号"
            style="width: 300px"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="filteredStudents" style="width: 100%" v-loading="loading">
        <el-table-column prop="student_id" label="学号" width="150" />
        <el-table-column prop="name" label="姓名" width="150" />
        <el-table-column label="专业">
          <template #default="{ row }">
            {{ row.major || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="deleteStudent(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { studentAPI } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const studentList = ref([])
const loading = ref(false)
const searchText = ref('')

const filteredStudents = computed(() => {
  if (!searchText.value) return studentList.value
  const search = searchText.value.toLowerCase()
  return studentList.value.filter(s => 
    s.name.toLowerCase().includes(search) || 
    s.student_id.toLowerCase().includes(search)
  )
})

const loadStudents = async () => {
  loading.value = true
  try {
    studentList.value = await studentAPI.getList()
  } catch (error) {
    ElMessage.error('加载学生列表失败')
  } finally {
    loading.value = false
  }
}

const deleteStudent = (row) => {
  ElMessageBox.confirm(
    `确定要删除学生"${row.name}"(${row.student_id})吗？这将同时删除其所有选课记录。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await studentAPI.delete(row.student_id)
      ElMessage.success('删除成功')
      loadStudents()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

onMounted(() => {
  loadStudents()
})
</script>

<style scoped>
.students-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
