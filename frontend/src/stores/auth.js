import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user  = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)

  function _persist(u, t) {
    user.value  = u
    token.value = t
    localStorage.setItem('user',  JSON.stringify(u))
    localStorage.setItem('token', t)
  }

  async function login(email, password) {
    const { data } = await authApi.login(email, password)
    _persist(data.user, data.token)
  }

  async function register(name, email, password) {
    const { data } = await authApi.register(name, email, password)
    _persist(data.user, data.token)
  }

  async function fetchMe() {
    const { data } = await authApi.me()
    user.value = data
    localStorage.setItem('user', JSON.stringify(data))
  }

  function logout() {
    user.value  = null
    token.value = null
    localStorage.removeItem('user')
    localStorage.removeItem('token')
  }

  return { user, token, isAuthenticated, login, register, fetchMe, logout }
})
