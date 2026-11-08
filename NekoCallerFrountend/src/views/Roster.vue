<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NCard, NDataTable, NIcon, NTag } from 'naive-ui'
import { PhArrowLeft, PhDownload, PhTrash } from '@phosphor-icons/vue'
import { classAPI, rosterAPI } from '@/api'
import { dialog, message } from '@/utils/feedback'
import { useBreakpoints } from '@/composables/useBreakpoints'

const route = useRoute()
const router = useRouter()
const { isCompact, isNarrow } = useBreakpoints()

const classId = route.params.classId
const className = ref('')
const roster = ref([])
const loading = ref(false)

async function loadRoster() {
  loading.value = true
  try {
    const [rosterData, classData] = await Promise.all([
      rosterAPI.getClassRoster(classId),
      classAPI.getDetail(classId)
    ])
    roster.value = rosterData || []
    className.value = classData?.class_name || ''
  } catch (error) {
    // 提示由拦截器统一弹
    console.error('加载花名册失败:', error)
  } finally {
    loading.value = false
  }
}

async function exportRoster() {
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
    message.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

function removeStudent(row) {
  dialog.warning({
    title: '移出班级',
    content: `确定把「${row.student_info.name}」从班级里移除吗？学生的积分记录会一起没掉。`,
    positiveText: '移除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await rosterAPI.removeStudent(row.enrollment_info.enrollment_id)
        message.success('移除成功')
        await loadRoster()
      } catch (error) {
        console.error('移除失败:', error)
      }
    }
  })
}

const columns = computed(() => {
  const base = [
    { title: '学号', key: 'student_id', width: isCompact.value ? 108 : 140, render: (row) => row.student_info.student_id },
    { title: '姓名', key: 'name', minWidth: isNarrow.value ? 76 : 100, render: (row) => row.student_info.name },
    { title: '专业', key: 'major', minWidth: 150, render: (row) => row.student_info.major || '—' },
    {
      title: '总积分',
      key: 'total_points',
      width: isCompact.value ? 78 : 100,
      sorter: (a, b) => a.enrollment_info.total_points - b.enrollment_info.total_points,
      render: (row) => row.enrollment_info.total_points
    },
    {
      title: '点名次数',
      key: 'call_count',
      width: isCompact.value ? 84 : 100,
      sorter: (a, b) => a.enrollment_info.call_count - b.enrollment_info.call_count,
      render: (row) => row.enrollment_info.call_count
    },
    {
      title: '转移权',
      key: 'transfer_rights',
      width: 88,
      hideOn: 'compact',
      render: (row) => {
        const rights = row.enrollment_info.transfer_rights || 0
        return rights > 0
          ? h(NTag, { size: 'small', type: 'success', bordered: false }, { default: () => rights })
          : h('span', { style: 'color:#8b95a7' }, '0')
      }
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
            title: '移除',
            onClick: () => removeStudent(row)
          },
          {
            icon: () => h(NIcon, null, { default: () => h(PhTrash) }),
            ...(isCompact.value ? {} : { default: () => '移除' })
          }
        )
    }
  ]

  // 窄屏先丢专业，再丢转移权，主要看学号姓名和积分
  return base.filter((column) => {
    if (column.hideOn === 'compact' && isCompact.value) return false
    if (column.key === 'major' && isCompact.value) return false
    return true
  })
})

onMounted(loadRoster)
</script>

<template>
  <div>
    <header class="page-head">
      <div class="head-row">
        <n-button quaternary size="small" @click="router.back()">
          <template #icon>
            <n-icon><PhArrowLeft /></n-icon>
          </template>
          返回
        </n-button>

        <div>
          <h2>{{ className || '班级花名册' }}</h2>
          <p>共 {{ roster.length }} 人</p>
        </div>

        <n-button secondary size="small" class="export" @click="exportRoster">
          <template #icon>
            <n-icon><PhDownload /></n-icon>
          </template>
          导出
        </n-button>
      </div>
    </header>

    <n-card size="small">
      <n-data-table
        :columns="columns"
        :data="roster"
        :loading="loading"
        :bordered="false"
        :row-key="(row) => row.enrollment_info.enrollment_id"
        :scroll-x="isNarrow ? 452 : isCompact ? 560 : 760"
        size="medium"
      />
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.head-row {
  display: flex;
  align-items: center;
  gap: 12px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
  }

  p {
    margin-top: 2px;
    font-size: 13px;
    color: #8b95a7;
  }
}

.export {
  margin-left: auto;
}
</style>
