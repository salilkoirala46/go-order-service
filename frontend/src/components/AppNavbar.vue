<template>
  <header class="bg-white border-b border-gray-200 px-6 py-3 flex items-center justify-between flex-shrink-0">
    <h1 class="text-lg font-semibold text-gray-800">{{ pageTitle }}</h1>
    <div class="flex items-center gap-3">
      <!-- Notification bell -->
      <RouterLink to="/notifications" class="relative p-2 text-gray-500 hover:text-gray-900 transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <span
          v-if="notifStore.unreadCount > 0"
          class="absolute top-1 right-1 bg-red-500 text-white text-xs w-4 h-4 flex items-center justify-center rounded-full font-bold"
        >
          {{ notifStore.unreadCount > 9 ? '9+' : notifStore.unreadCount }}
        </span>
      </RouterLink>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useNotificationsStore } from '@/stores/notifications'
import { onMounted } from 'vue'

const route      = useRoute()
const notifStore = useNotificationsStore()

const titles = {
  '/dashboard':     'Dashboard',
  '/orders':        'Orders',
  '/orders/new':    'New Order',
  '/notifications': 'Notifications',
  '/profile':       'Profile'
}

const pageTitle = computed(() => {
  if (route.name === 'order-detail') return 'Order Details'
  return titles[route.path] ?? 'Order Services'
})

onMounted(() => notifStore.fetchNotifications())
</script>
