<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NButton, NCard, NDataTable, NIcon, NInput } from 'naive-ui'
import { PhMagnifyingGlass, PhTrash } from '@phosphor-icons/vue'
import { studentAPI } from '@/api'
import { dialog, message } from '@/utils/feedback'
import { useBreakpoints } from '@/composables/useBreakpoints'

const { isCompact, isNarrow } = useBreakpoints()

const studentList = ref([])
const loading = ref(false)
const searchText = ref('')

const filteredStudents = computed(() => {
  const keyword = searchText.value.trim().toLowerCase()
  if (!keyword) return studentList.value

  return studentList.value.filter(
    (item) =>
      item.name?.toLowerCase().includes(keyword) ||
      item.student_id?.toLowerCase().includes(keyword)
  )
})

async function loadStudents() {
  loading.value = true
  try {
    studentList.value = await studentAPI.getList()
  } catch (error) {
    // 提示由拦截器统一弹
    console.error('加载学生列表失败:', error)
  } finally {
    loading.value = false
  }
}

function deleteStudent(row) {
  dialog.warning({
    title: '删除学生',
    content: `确定删除「${row.name}」吗？该学生在所有班级的记录会一起清掉。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await studentAPI.delete(row.student_id)
        message.success('删除成功')
        await loadStudents()
      } catch (error) {
        console.error('删除学生失败:', error)
      }
    }
  })
}

const columns = computed(() =>
  [
  { title: '学号', key: 'student_id', width: isCompact.value ? 112 : 150 },
  { title: '姓名', key: 'name', minWidth: isNarrow.value ? 80 : 120 },
  {
    title: '专业',
    key: 'major',
    minWidth: 160,
    render: (row) => row.major || '—',
    hideOn: 'compact'
  },
  {
    title: '操作',
    key: 'actions',
    width: isCompact.value ? 80 : 100,
    render: (row) =>
      h(
        NButton,
        {
          size: 'tiny',
          type: 'error',
          quaternary: true,
          circle: isCompact.value,
          title: '删除',
          onClick: () => deleteStudent(row)
        },
        {
          icon: () => h(NIcon, null, { default: () => h(PhTrash) }),
          ...(isCompact.value ? {} : { default: () => '删除' })
        }
      )
  }
].filter((column) => !(column.hideOn === 'compact' && isCompact.value))
)

onMounted(loadStudents)
</script>

<template>
  <div>
    <header class="page-head">
      <h2>学生管理</h2>
      <p>所有导入过的学生，按姓名或学号搜</p>
    </header>

    <n-card size="small" :title="`学生列表 · ${filteredStudents.length} 人`">
      <template #header-extra>
        <n-input
          v-model:value="searchText"
          placeholder="搜索姓名或学号"
          clearable
          class="search"
        >
          <template #prefix>
            <n-icon><PhMagnifyingGlass /></n-icon>
          </template>
        </n-input>
      </template>

      <n-data-table
        :columns="columns"
        :data="filteredStudents"
        :loading="loading"
        :bordered="false"
        :row-key="(row) => row.student_id"
        :scroll-x="isNarrow ? 420 : isCompact ? 480 : 640"
        size="medium"
      />
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.search {
  width: 240px;
}

@media (max-width: 720px) {
  .search {
    width: 140px;
  }
}
</style>
