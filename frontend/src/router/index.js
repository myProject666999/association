import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', requiresAuth: true }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/users/List.vue'),
        meta: { title: '用户管理', requiresAuth: true, roles: ['university_admin', 'dept_admin'] }
      },
      {
        path: 'clubs',
        name: 'Clubs',
        component: () => import('@/views/clubs/List.vue'),
        meta: { title: '社团列表', requiresAuth: true }
      },
      {
        path: 'clubs/my',
        name: 'MyClubs',
        component: () => import('@/views/clubs/MyClubs.vue'),
        meta: { title: '我的社团', requiresAuth: true }
      },
      {
        path: 'clubs/:id',
        name: 'ClubDetail',
        component: () => import('@/views/clubs/Detail.vue'),
        meta: { title: '社团详情', requiresAuth: true }
      },
      {
        path: 'activities',
        name: 'Activities',
        component: () => import('@/views/activities/List.vue'),
        meta: { title: '活动列表', requiresAuth: true }
      },
      {
        path: 'activities/my',
        name: 'MyActivities',
        component: () => import('@/views/activities/MyActivities.vue'),
        meta: { title: '我的活动', requiresAuth: true }
      },
      {
        path: 'activities/organized',
        name: 'OrganizedActivities',
        component: () => import('@/views/activities/Organized.vue'),
        meta: { title: '我组织的活动', requiresAuth: true }
      },
      {
        path: 'activities/:id',
        name: 'ActivityDetail',
        component: () => import('@/views/activities/Detail.vue'),
        meta: { title: '活动详情', requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/Index.vue'),
        meta: { title: '个人中心', requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 大学生社团管理系统` : '大学生社团管理系统'

  const userStore = useUserStore()
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth) {
    if (!token) {
      next({ name: 'Login' })
      return
    }

    if (!userStore.user) {
      try {
        await userStore.fetchUserInfo()
      } catch (error) {
        localStorage.removeItem('token')
        next({ name: 'Login' })
        return
      }
    }

    if (to.meta.roles && to.meta.roles.length > 0) {
      if (!to.meta.roles.includes(userStore.user.role)) {
        next({ name: 'Dashboard' })
        return
      }
    }
  }

  if ((to.name === 'Login' || to.name === 'Register') && token) {
    next({ name: 'Dashboard' })
    return
  }

  next()
})

export default router
