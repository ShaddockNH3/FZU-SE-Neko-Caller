import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/classes',
      name: 'classes',
      component: () => import('../views/Classes.vue')
    },
    {
      path: '/students',
      name: 'students',
      component: () => import('../views/Students.vue')
    },
    {
      path: '/roll-call',
      name: 'roll-call',
      component: () => import('../views/RollCall.vue')
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: () => import('../views/Leaderboard.vue')
    },
    {
      path: '/roster/:classId',
      name: 'roster',
      component: () => import('../views/Roster.vue')
    }
  ]
})

export default router
