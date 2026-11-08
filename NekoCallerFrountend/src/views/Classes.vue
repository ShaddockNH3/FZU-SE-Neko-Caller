<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NModal,
  NTabPane,
  NTabs,
  NUpload,
  NUploadDragger
} from 'naive-ui'
import {
  PhDownload,
  PhListDashes,
  PhTrash,
  PhTrophy,
  PhUploadSimple
} from '@phosphor-icons/vue'
import { classAPI, importAPI } from '@/api'
import { dialog, message } from '@/utils/feedback'
import { useBreakpoints } from '@/composables/useBreakpoints'

const router = useRouter()
const { isCompact } = useBreakpoints()

const classList = ref([])
const loading = ref(false)
const showImport = ref(false)
const importType = ref('excel')
const importing = ref(false)

const excelForm = ref({ className: '', file: null })
const jsonForm = ref({ className: '', students: '' })

async function loadClasses() {
  loading.value = true
  try {
    classList.value = await classAPI.getList()
  } catch (error) {
    // 提示由拦截器统一弹
    console.error('加载班级列表失败:', error)
  } finally {
    loading.value = false
  }
}

function onFileChange({ fileList }) {
  excelForm.value.file = fileList[fileList.length - 1]?.file || null
}

async function handleImport() {
  if (importType.value === 'excel') {
    if (!excelForm.value.className) return message.warning('请输入班级名称')
    if (!excelForm.value.file) return message.warning('请选择 Excel 文件')

    importing.value = true
    try {
      const formData = new FormData()
      formData.append('file', excelForm.value.file)
      formData.append('class_name', excelForm.value.className)

      const res = await importAPI.importExcel(formData)
      message.success(res?.message || '导入成功')
      showImport.value = false
      excelForm.value = { className: '', file: null }
      await loadClasses()
    } catch (error) {
      console.error('Excel 导入失败:', error)
    } finally {
      importing.value = false
    }
    return
  }

  if (!jsonForm.value.className || !jsonForm.value.students) {
    return message.warning('班级名称和学生数据都要填')
  }

  importing.value = true
  try {
    let students
    try {
      students = JSON.parse(jsonForm.value.students)
    } catch {
      return message.error('学生数据不是合法 JSON')
    }

    await importAPI.importJSON({ class_name: jsonForm.value.className, students })
    message.success('导入成功')
    showImport.value = false
    jsonForm.value = { className: '', students: '' }
    await loadClasses()
  } catch (error) {
    console.error('JSON 导入失败:', error)
  } finally {
    importing.value = false
  }
}

function viewRoster(classId) {
  router.push(`/roster/${classId}`)
}

function viewLeaderboard(classId) {
  router.push({ path: '/leaderboard', query: { classId } })
}

async function exportExcel(row) {
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
    message.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

function deleteClass(row) {
  dialog.warning({
    title: '删除班级',
    content: `确定删除「${row.class_name}」吗？班级内的点名记录和积分会一起没掉。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await classAPI.delete(row.class_id)
        message.success('删除成功')
        await loadClasses()
      } catch (error) {
        console.error('删除失败:', error)
      }
    }
  })
}

/** 窄屏只留图标，文案进 title。 */
function actionButton({ label, icon, type, secondary, onClick }) {
  const iconOnly = isCompact.value

  return h(
    NButton,
    {
      size: 'tiny',
      type,
      secondary: secondary ?? true,
      circle: iconOnly,
      title: label,
      style: 'margin-right:6px',
      onClick
    },
    {
      icon: () => h(NIcon, null, { default: () => h(icon) }),
      ...(iconOnly ? {} : { default: () => label })
    }
  )
}

const columns = computed(() => {
  const base = [
    { title: '班级名称', key: 'class_name', minWidth: 180 },
    {
      title: '学生人数',
      key: 'student_count',
      width: 110,
      render: (row) => row.student_ids?.length || 0
    },
    {
      title: '班级 ID',
      key: 'class_id',
      width: 300,
      className: 'nowrap',
      render: (row) => h('span', { style: 'font-family:ui-monospace,monospace;font-size:12px;color:#8b95a7' }, row.class_id)
    },
    {
      title: '操作',
      key: 'actions',
      width: isCompact.value ? 170 : 340,
      render: (row) =>
        h('div', { style: 'display:flex;flex-wrap:wrap;row-gap:6px' }, [
          actionButton({
            label: '花名册',
            icon: PhListDashes,
            onClick: () => viewRoster(row.class_id)
          }),
          actionButton({
            label: '排行榜',
            icon: PhTrophy,
            type: 'success',
            onClick: () => viewLeaderboard(row.class_id)
          }),
          actionButton({
            label: '导出',
            icon: PhDownload,
            type: 'warning',
            onClick: () => exportExcel(row)
          }),
          actionButton({
            label: '删除',
            icon: PhTrash,
            type: 'error',
            onClick: () => deleteClass(row)
          })
        ])
    }
  ]

  // 窄屏先把 ID 列拿掉，那一列又长又没人看
  return isCompact.value ? base.filter((column) => column.key !== 'class_id') : base
})

onMounted(loadClasses)
</script>

<template>
  <div>
    <header class="page-head">
      <h2>班级管理</h2>
      <p>导入班级名单，进花名册看每个学生</p>
    </header>

    <n-card size="small" title="班级列表">
      <template #header-extra>
        <n-button type="primary" size="small" @click="showImport = true">
          <template #icon>
            <n-icon><PhUploadSimple /></n-icon>
          </template>
          导入班级
        </n-button>
      </template>

      <n-data-table
        :columns="columns"
        :data="classList"
        :loading="loading"
        :bordered="false"
        :row-key="(row) => row.class_id"
        :scroll-x="isCompact ? 720 : 1060"
        size="medium"
      />
    </n-card>

    <n-modal
      v-model:show="showImport"
      preset="card"
      title="导入班级数据"
      style="width: 560px; max-width: 92vw"
    >
      <n-tabs v-model:value="importType" type="segment">
        <n-tab-pane name="excel" tab="Excel 导入">
          <n-form label-placement="left" :label-width="80" class="form">
            <n-form-item label="班级名称">
              <n-input v-model:value="excelForm.className" placeholder="例如：2025 软件工程 K 班" />
            </n-form-item>

            <n-form-item label="名单文件">
              <n-upload
                :default-upload="false"
                :max="1"
                accept=".xlsx,.xls"
                @change="onFileChange"
              >
                <n-upload-dragger>
                  <div class="dragger">
                    <n-icon :size="28" color="#8b95a7"><PhUploadSimple /></n-icon>
                    <p>把 xlsx 拖进来，或者点击选择</p>
                    <span>列顺序：学号 · 姓名 · 专业（可选）</span>
                  </div>
                </n-upload-dragger>
              </n-upload>
            </n-form-item>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="json" tab="JSON 导入">
          <n-form label-placement="left" :label-width="80" class="form">
            <n-form-item label="班级名称">
              <n-input v-model:value="jsonForm.className" placeholder="例如：2025 软件工程 K 班" />
            </n-form-item>

            <n-form-item label="学生数据">
              <n-input
                v-model:value="jsonForm.students"
                type="textarea"
                :rows="10"
                placeholder='[{"student_id":"83210001","name":"张三","major":"软件工程"}]'
              />
            </n-form-item>
          </n-form>
        </n-tab-pane>
      </n-tabs>

      <template #footer>
        <div class="modal-footer">
          <n-button quaternary @click="showImport = false">取消</n-button>
          <n-button type="primary" :loading="importing" @click="handleImport">确定导入</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<style lang="scss" scoped>
.form {
  margin-top: 12px;
}

.dragger {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 0;
  color: #4b5563;

  p {
    font-size: 13px;
  }

  span {
    font-size: 12px;
    color: #8b95a7;
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
