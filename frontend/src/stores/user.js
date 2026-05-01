import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  async function login(username, password) {
    const response = await api.post('/api/login', { username, password })
    if (response.code === 200) {
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', response.data.token)
      return true
    }
    throw new Error(response.message)
  }

  async function register(userData) {
    const response = await api.post('/api/register', userData)
    return response.code === 200
  }

  async function fetchUserInfo() {
    const response = await api.get('/api/user/me')
    if (response.code === 200) {
      user.value = response.data
      return true
    }
    throw new Error(response.message)
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }

  function isAdmin() {
    return user.value && (user.value.role === 'university_admin' || user.value.role === 'dept_admin')
  }

  function isUniversityAdmin() {
    return user.value && user.value.role === 'university_admin'
  }

  function isDeptAdmin() {
    return user.value && user.value.role === 'dept_admin'
  }

  return {
    user,
    token,
    login,
    register,
    fetchUserInfo,
    logout,
    isAdmin,
    isUniversityAdmin,
    isDeptAdmin
  }
})
