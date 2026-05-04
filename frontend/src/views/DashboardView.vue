<template>
  <div class="space-y-6">
    <!-- Welcome -->
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Welcome back, {{ auth.user?.name?.split(' ')[0] }} 👋</h2>
      <p class="text-gray-500 text-sm mt-1">Here's what's happening with your orders.</p>
    </div>

    <!-- Stats -->
    <LoadingSpinner v-if="ordersStore.loading" />
    <div v-else class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in stats" :key="stat.label" class="card">
        <p class="text-sm text-gray-500 mb-1">{{ stat.label }}</p>
        <p class="text-3xl font-bold" :class="stat.color">{{ stat.value }}</p>
      </div>
    </div>

    <!-- Recent orders -->
    <div class="card">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-semibold text-gray-900">Recent Orders</h3>
        <RouterLink to="/orders" class="text-sm text-blue-600 hover:underline">View all</RouterLink>
      </div>

      <EmptyState
        v-if="!ordersStore.loading && recent.length === 0"
        title="No orders yet"
        message="Create your first order to get started."
      >
        <RouterLink to="/orders/new" class="btn-primary text-sm">New Order</RouterLink>
      </EmptyState>

      <div v-else class="divide-y divide-gray-100">
        <div
          v-for="order in recent"
          :key="order.id"
          class="flex items-center justify-between py-3 cursor-pointer hover:bg-gray-50 -mx-6 px-6 transition-colors"
          @click="router.push(`/orders/${order.id}`)"
        >
          <div>
            <p class="font-medium text-gray-900 text-sm">{{ order.product }}</p>
            <p class="text-xs text-gray-400">{{ formatDate(order.created_at) }} · qty {{ order.quantity }}</p>
          </div>
          <div class="flex items-center gap-3">
            <span class="font-semibold text-sm">${{ order.total.toFixed(2) }}</span>
            <StatusBadge :status="order.status" />
          </div>
        </div>
      </div>
    </div>

    <!-- Notifications preview -->
    <div class="card">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-semibold text-gray-900">
          Notifications
          <span v-if="notifStore.unreadCount > 0" class="ml-2 text-xs bg-red-500 text-white px-1.5 py-0.5 rounded-full">
            {{ notifStore.unreadCount }} new
          </span>
        </h3>
        <RouterLink to="/notifications" class="text-sm text-blue-600 hover:underline">View all</RouterLink>
      </div>

      <EmptyState
        v-if="!notifStore.loading && notifStore.notifications.length === 0"
        title="No notifications"
        message="You're all caught up."
      />

      <div v-else class="space-y-2">
        <div
          v-for="n in recentNotifications"
          :key="n.id"
          :class="['flex gap-3 p-3 rounded-lg', n.read ? 'bg-white' : 'bg-blue-50']"
        >
          <div :class="['w-2 h-2 rounded-full mt-2 flex-shrink-0', n.read ? 'bg-gray-300' : 'bg-blue-500']" />
          <div>
            <p class="text-sm text-gray-800">{{ n.message }}</p>
            <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(n.created_at) }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useOrdersStore } from '@/stores/orders'
import { useNotificationsStore } from '@/stores/notifications'
import StatusBadge from '@/components/StatusBadge.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'

const auth         = useAuthStore()
const ordersStore  = useOrdersStore()
const notifStore   = useNotificationsStore()
const router       = useRouter()

const recent = computed(() => ordersStore.orders.slice(0, 5))
const recentNotifications = computed(() => notifStore.notifications.slice(0, 4))

const stats = computed(() => {
  const o = ordersStore.orders
  return [
    { label: 'Total Orders',    value: o.length,                                          color: 'text-gray-900' },
    { label: 'Pending',         value: o.filter(x => x.status === 'pending').length,      color: 'text-yellow-600' },
    { label: 'Delivered',       value: o.filter(x => x.status === 'delivered').length,    color: 'text-green-600' },
    { label: 'Total Spent',     value: '$' + o.reduce((s, x) => s + x.total, 0).toFixed(2), color: 'text-blue-600' }
  ]
})

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(async () => {
  await Promise.all([ordersStore.fetchOrders(), notifStore.fetchNotifications()])
})
</script>
