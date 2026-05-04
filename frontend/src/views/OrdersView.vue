<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">Orders</h2>
        <p class="text-sm text-gray-500">{{ store.orders.length }} total orders</p>
      </div>
      <RouterLink to="/orders/new" class="btn-primary gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Order
      </RouterLink>
    </div>

    <!-- Filter tabs -->
    <div class="flex gap-2 border-b border-gray-200">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="activeTab = tab.value"
        :class="[
          'px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors',
          activeTab === tab.value
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700'
        ]"
      >
        {{ tab.label }}
        <span class="ml-1.5 text-xs bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded-full">{{ tab.count }}</span>
      </button>
    </div>

    <LoadingSpinner v-if="store.loading" />

    <EmptyState
      v-else-if="filtered.length === 0"
      title="No orders found"
      message="Try a different filter or create a new order."
    >
      <RouterLink to="/orders/new" class="btn-primary text-sm">New Order</RouterLink>
    </EmptyState>

    <!-- Table -->
    <div v-else class="card p-0 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b border-gray-100">
          <tr>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Order</th>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Product</th>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden sm:table-cell">Qty</th>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden md:table-cell">Total</th>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Status</th>
            <th class="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden lg:table-cell">Date</th>
            <th class="px-6 py-3"></th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr
            v-for="order in filtered"
            :key="order.id"
            class="hover:bg-gray-50 transition-colors cursor-pointer"
            @click="router.push(`/orders/${order.id}`)"
          >
            <td class="px-6 py-4 font-mono text-xs text-gray-500">#{{ order.id }}</td>
            <td class="px-6 py-4 font-medium text-gray-900">{{ order.product }}</td>
            <td class="px-6 py-4 text-gray-600 hidden sm:table-cell">{{ order.quantity }}</td>
            <td class="px-6 py-4 font-semibold hidden md:table-cell">${{ order.total.toFixed(2) }}</td>
            <td class="px-6 py-4"><StatusBadge :status="order.status" /></td>
            <td class="px-6 py-4 text-gray-400 hidden lg:table-cell">{{ formatDate(order.created_at) }}</td>
            <td class="px-6 py-4 text-right">
              <svg class="w-4 h-4 text-gray-400 ml-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import StatusBadge from '@/components/StatusBadge.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'

const store     = useOrdersStore()
const router    = useRouter()
const activeTab = ref('all')

const tabs = computed(() => [
  { value: 'all',       label: 'All',       count: store.orders.length },
  { value: 'pending',   label: 'Pending',   count: store.orders.filter(o => o.status === 'pending').length },
  { value: 'confirmed', label: 'Confirmed', count: store.orders.filter(o => o.status === 'confirmed').length },
  { value: 'shipped',   label: 'Shipped',   count: store.orders.filter(o => o.status === 'shipped').length },
  { value: 'delivered', label: 'Delivered', count: store.orders.filter(o => o.status === 'delivered').length },
  { value: 'cancelled', label: 'Cancelled', count: store.orders.filter(o => o.status === 'cancelled').length }
])

const filtered = computed(() =>
  activeTab.value === 'all'
    ? store.orders
    : store.orders.filter(o => o.status === activeTab.value)
)

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(() => store.fetchOrders())
</script>
