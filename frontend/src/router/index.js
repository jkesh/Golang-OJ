import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/problems',
    name: 'Problems',
    component: () => import('../views/Problems.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/submissions',
    name: 'Submission',
    component: () => import('../views/Submission.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/submission/:id',
    name: 'SubmissionDetail',
    component: () => import('../views/SubmissionDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
      path: '/problem/:id',
      name: 'ProblemDetail',
      component: () => import('../views/ProblemDetail.vue'),
      meta: { requiresAuth: true }
  }

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫，检查用户是否已登录
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router