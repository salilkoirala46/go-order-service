import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ordersApi } from '@/api/orders'

export const useOrdersStore = defineStore('orders', () => {
  const orders  = ref([])
  const current = ref(null)
  const loading = ref(false)
  const error   = ref(null)

  async function fetchOrders() {
    loading.value = true
    error.value   = null
    try {
      const { data } = await ordersApi.list()
      orders.value = data ?? []
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Failed to load orders'
    } finally {
      loading.value = false
    }
  }

  async function fetchOrder(id) {
    loading.value = true
    error.value   = null
    try {
      const { data } = await ordersApi.get(id)
      current.value = data
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Failed to load order'
    } finally {
      loading.value = false
    }
  }

  async function createOrder(product, quantity, price) {
    const { data } = await ordersApi.create(product, quantity, price)
    orders.value.unshift(data)
    return data
  }

  async function updateStatus(id, status) {
    const { data } = await ordersApi.updateStatus(id, status)
    const idx = orders.value.findIndex(o => o.id === id)
    if (idx !== -1) orders.value[idx] = data
    if (current.value?.id === id) current.value = data
    return data
  }

  return { orders, current, loading, error, fetchOrders, fetchOrder, createOrder, updateStatus }
})
