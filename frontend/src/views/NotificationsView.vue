<template>
  <div class="max-w-2xl mx-auto space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">Notifications</h2>
        <p class="text-sm text-gray-500">
          {{ store.unreadCount > 0 ? `${store.unreadCount} unread` : 'All caught up' }}
        </p>
      </div>
      <button
        v-if="store.unreadCount > 0"
        @click="store.markAllRead()"
        class="btn-secondary text-sm"
      >
        Mark all as read
      </button>
    </div>

    <LoadingSpinner v-if="store.loading" />

    <EmptyState
      v-else-if="store.notifications.length === 0"
      title="No notifications yet"
      message="Notifications appear here when you register, place orders, or update order statuses."
    />

    <div v-else class="space-y-2">
      <div
        v-for="n in store.notifications"
        :key="n.id"
        :class="['card p-4 flex gap-4 cursor-pointer hover:shadow-md transition-shadow', !n.read ? 'border-blue-200 bg-blue-50/40' : '']"
        @click="!n.read && store.markRead(n.id)"
      >
        <!-- Type icon -->
        <div :class="['w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0', typeStyle(n.type).bg]">
          <span class="text-lg">{{ typeStyle(n.type).emoji }}</span>
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-start justify-between gap-2">
            <p :class="['text-sm', n.read ? 'text-gray-600' : 'text-gray-900 font-medium']">
              {{ n.message }}
            </p>
            <div class="flex-shrink-0 flex items-center gap-2">
              <span v-if="!n.read" class="w-2 h-2 bg-blue-500 rounded-full" />
            </div>
          </div>
          <div class="flex items-center gap-2 mt-1">
            <span :class="['text-xs font-medium px-1.5 py-0.5 rounded', typeStyle(n.type).badge]">
              {{ formatType(n.type) }}
            </span>
            <span class="text-xs text-gray-400">{{ formatDate(n.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'

const store = useNotificationsStore()

const typeStyles = {
  WELCOME:             { bg: 'bg-green-100',  badge: 'bg-green-100 text-green-700',  emoji: '👋' },
  ORDER_CREATED:       { bg: 'bg-blue-100',   badge: 'bg-blue-100 text-blue-700',    emoji: '📦' },
  ORDER_STATUS_CHANGE: { bg: 'bg-purple-100', badge: 'bg-purple-100 text-purple-700', emoji: '🔄' }
}

function typeStyle(type) {
  return typeStyles[type] ?? { bg: 'bg-gray-100', badge: 'bg-gray-100 text-gray-700', emoji: '🔔' }
}

function formatType(type) {
  return { WELCOME: 'Welcome', ORDER_CREATED: 'New Order', ORDER_STATUS_CHANGE: 'Status Update' }[type] ?? type
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => store.fetchNotifications())
</script>
