import { createRouter, createWebHistory } from 'vue-router'
import {
  PhGraduationCap,
  PhHouse,
  PhMegaphone,
  PhTrophy,
  PhUsers
} from '@phosphor-icons/vue'

// 菜单从这份路由生成（见 App.vue），加页面只加一条：meta.title 决定文案，meta.icon 决定图标。
const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue'),
    meta: { title: '首页', icon: PhHouse }
  },
  {
    path: '/classes',
    name: 'classes',
    component: () => import('../views/Classes.vue'),
    meta: { title: '班级管理', icon: PhGraduationCap }
  },
  {
    path: '/students',
    name: 'students',
    component: () => import('../views/Students.vue'),
    meta: { title: '学生管理', icon: PhUsers }
  },
  {
    path: '/roll-call',
    name: 'roll-call',
    component: () => import('../views/RollCall.vue'),
    meta: { title: '点名', icon: PhMegaphone }
  },
  {
    path: '/leaderboard',
    name: 'leaderboard',
    component: () => import('../views/Leaderboard.vue'),
    meta: { title: '排行榜', icon: PhTrophy }
  },
  {
    path: '/roster/:classId',
    name: 'roster',
    component: () => import('../views/Roster.vue'),
    meta: { title: '花名册', parent: '/classes' }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · Neko Caller` : 'Neko Caller'
})

export default router
