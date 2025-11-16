<template>
  <div class="classes-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>班级管理</span>
          <div>
            <el-button type="primary" @click="showImportDialog = true">
              <el-icon><Upload /></el-icon>
              导入班级
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="classList" style="width: 100%" v-loading="loading">
        <el-table-column prop="class_name" label="班级名称" />
        <el-table-column label="学生人数">
          <template #default="{ row }">
            {{ row.student_ids?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="class_id" label="班级ID" width="280" />
        <el-table-column label="操作" width="350">
          <template #default="{ row }">
            <el-button size="small" @click="viewRoster(row.class_id)">
              花名册
            </el-button>
            <el-button size="small" type="success" @click="viewLeaderboard(row.class_id)">
              排行榜
            </el-button>
            <el-button size="small" type="warning" @click="exportExcel(row)">
              导出
            </el-button>
            <el-button size="small" type="danger" @click="deleteClass(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 导入对话框 -->
    <el-dialog v-model="showImportDialog" title="导入班级数据" width="500px">
      <el-tabs v-model="importType">
        <el-tab-pane label="Excel导入" name="excel">
          <el-form :model="excelForm" label-width="100px">
            <el-form-item label="班级名称" required>
              <el-input v-model="excelForm.className" placeholder="请输入班级名称" />
            </el-form-item>
            <el-form-item label="Excel文件" required>
              <el-upload
                ref="uploadRef"
                :auto-upload="false"
                :limit="1"
                :on-change="handleFileChange"
                accept=".xlsx,.xls"
                drag
              >
                <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
                <div class="el-upload__text">
                  拖拽文件到这里或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    只支持 xlsx/xls 文件，格式: 学号 | 姓名 | 专业(可选)
                  </div>
                </template>
              </el-upload>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="JSON导入" name="json">
          <el-form :model="jsonForm" label-width="100px">
            <el-form-item label="班级名称" required>
              <el-input v-model="jsonForm.className" placeholder="请输入班级名称" />
            </el-form-item>
            <el-form-item label="学生数据" required>
              <el-input
                v-model="jsonForm.students"
                type="textarea"
                :rows="10"
                placeholder='[{"student_id":"001","name":"张三","major":"软件工程"}]'
              />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importing">
          确定导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { classAPI, importAPI } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const classList = ref([])
const loading = ref(false)
const showImportDialog = ref(false)
const importType = ref('excel')
const importing = ref(false)

const excelForm = ref({
  className: '',
  file: null
})

const jsonForm = ref({
  className: '',
  students: ''
})

const loadClasses = async () => {
  loading.value = true
  try {
    classList.value = await classAPI.getList()
  } catch (error) {
    ElMessage.error('加载班级列表失败')
  } finally {
    loading.value = false
  }
}

const handleFileChange = (file) => {
  excelForm.value.file = file.raw
}

const handleImport = async () => {
  if (importType.value === 'excel') {
    if (!excelForm.value.className) {
      ElMessage.warning('请输入班级名称')
      return
    }
    if (!excelForm.value.file) {
      ElMessage.warning('请选择Excel文件')
      return
    }
    
    importing.value = true
    try {
      const formData = new FormData()
      formData.append('file', excelForm.value.file)
      formData.append('class_name', excelForm.value.className)
      
      const res = await importAPI.importExcel(formData)
      ElMessage.success(res.message || '导入成功')
      showImportDialog.value = false
      excelForm.value = { className: '', file: null }
      loadClasses()
    } catch (error) {
      ElMessage.error('导入失败')
    } finally {
      importing.value = false
    }
  } else {
    if (!jsonForm.value.className || !jsonForm.value.students) {
      ElMessage.warning('请填写完整信息')
      return
    }
    
    importing.value = true
    try {
      const students = JSON.parse(jsonForm.value.students)
      await importAPI.importJSON({
        class_name: jsonForm.value.className,
        students
      })
      ElMessage.success('导入成功')
      showImportDialog.value = false
      jsonForm.value = { className: '', students: '' }
      loadClasses()
    } catch (error) {
      ElMessage.error('导入失败，请检查JSON格式')
    } finally {
      importing.value = false
    }
  }
}

const viewRoster = (classId) => {
  router.push(`/roster/${classId}`)
}

const viewLeaderboard = (classId) => {
  router.push({ path: '/leaderboard', query: { classId } })
}

const exportExcel = async (row) => {
  try {
    const response = await classAPI.exportExcel(row.class_id)
    const blob = new Blob([response.data], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${row.class_name}-积分详单.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

const deleteClass = (row) => {
  ElMessageBox.confirm(
    `确定要删除班级"${row.class_name}"吗？这将同时删除所有相关数据。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await classAPI.delete(row.class_id)
      ElMessage.success('删除成功')
      loadClasses()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

onMounted(() => {
  loadClasses()
})
</script>

<style scoped>
.classes-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
