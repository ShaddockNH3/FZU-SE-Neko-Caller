import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/v1',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    // 如果返回的是文件流
    if (response.config.responseType === 'blob') {
      return response
    }
    // 正常数据返回
    return res
  },
  error => {
    ElMessage.error(error.response?.data?.message || '请求失败')
    return Promise.reject(error)
  }
)

// 班级相关API
export const classAPI = {
  // 获取班级列表
  getList() {
    return request.get('/classes')
  },
  // 获取班级详情
  getDetail(classId) {
    return request.get(`/classes/${classId}`)
  },
  // 删除班级
  delete(classId) {
    return request.delete(`/classes/${classId}`)
  },
  // 获取排行榜
  getLeaderboard(classId, top = 10) {
    return request.get(`/classes/${classId}/leaderboard`, { params: { top } })
  },
  // 获取统计信息
  getStats(classId) {
    return request.get(`/classes/${classId}/stats`)
  },
  // 导出Excel
  exportExcel(classId) {
    return request.get(`/classes/${classId}/export`, { responseType: 'blob' })
  }
}

// 学生相关API
export const studentAPI = {
  // 获取学生列表
  getList() {
    return request.get('/students')
  },
  // 获取学生详情
  getDetail(studentId) {
    return request.get(`/students/${studentId}`)
  },
  // 删除学生
  delete(studentId) {
    return request.delete(`/students/${studentId}`)
  }
}

// 花名册相关API
export const rosterAPI = {
  // 获取班级花名册
  getClassRoster(classId) {
    return request.get('/roster', { params: { class_id: classId } })
  },
  // 移除学生
  removeStudent(enrollmentId) {
    return request.delete(`/enrollments/${enrollmentId}`)
  }
}

// 点名相关API
export const rollCallAPI = {
  // 发起点名
  call(data) {
    return request.post('/roll-calls', data)
  },
  // 结算点名
  solve(data) {
    return request.post('/roll-calls/solve', data)
  },
  // 重置点名状态
  reset(classId) {
    return request.post('/roll-calls/reset', { class_id: classId })
  }
}

// 导入相关API
export const importAPI = {
  // JSON导入
  importJSON(data) {
    return request.post('/import/class-data', data)
  },
  // Excel导入
  importExcel(formData) {
    return request.post('/import/excel', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}

export default request
