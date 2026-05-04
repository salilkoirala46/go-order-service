import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { notificationsApi } from '@/api/notifications'

export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref([])
  const loading       = ref(false)

  const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)

  async function fetchNotifications() {
    loading.value = true
    try {
      const { data } = await notificationsApi.list()
      notifications.value = data ?? []
    } finally {
      loading.value = false
    }
  }

  async function markRead(id) {
    await notificationsApi.markRead(id)
    const n = notifications.value.find(n => n.id === id)
    if (n) n.read = true
  }

  async function markAllRead() {
    const unread = notifications.value.filter(n => !n.read)
    await Promise.all(unread.map(n => notificationsApi.markRead(n.id)))
    unread.forEach(n => { n.read = true })
  }

  return { notifications, loading, unreadCount, fetchNotifications, markRead, markAllRead }
})
