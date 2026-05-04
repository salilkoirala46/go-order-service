<template>
  <div class="max-w-lg mx-auto space-y-5">
    <div>
      <h2 class="text-xl font-bold text-gray-900">Profile</h2>
      <p class="text-sm text-gray-500">Your account information</p>
    </div>

    <div class="card">
      <!-- Avatar -->
      <div class="flex items-center gap-4 mb-6 pb-6 border-b border-gray-100">
        <div class="w-16 h-16 bg-blue-600 rounded-full flex items-center justify-center text-2xl font-bold text-white">
          {{ initials }}
        </div>
        <div>
          <h3 class="text-lg font-bold text-gray-900">{{ auth.user?.name }}</h3>
          <p class="text-sm text-gray-500">{{ auth.user?.email }}</p>
          <p class="text-xs text-gray-400 mt-0.5">Member since {{ joinDate }}</p>
        </div>
      </div>

      <dl class="space-y-4">
        <div class="flex justify-between text-sm">
          <dt class="text-gray-500">Full name</dt>
          <dd class="font-medium text-gray-900">{{ auth.user?.name }}</dd>
        </div>
        <div class="flex justify-between text-sm">
          <dt class="text-gray-500">Email address</dt>
          <dd class="font-medium text-gray-900">{{ auth.user?.email }}</dd>
        </div>
        <div class="flex justify-between text-sm">
          <dt class="text-gray-500">User ID</dt>
          <dd class="font-mono text-gray-600">#{{ auth.user?.id }}</dd>
        </div>
      </dl>
    </div>

    <!-- Order stats -->
    <div class="card">
      <h4 class="font-semibold text-gray-900 mb-4">Order Summary</h4>
      <div class="grid grid-cols-2 gap-4">
        <div v-for="stat in stats" :key="stat.label" class="bg-gray-50 rounded-lg p-3 text-center">
          <p class="text-2xl font-bold" :class="stat.color">{{ stat.value }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ stat.label }}</p>
        </div>
      </div>
    </div>

    <button @click="logout" class="btn-danger w-full">Sign out</button>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useOrdersStore } from '@/stores/orders'

const auth         = useAuthStore()
const ordersStore  = useOrdersStore()
const router       = useRouter()

const initials = computed(() =>
  auth.user?.name?.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2) ?? '?'
)

const joinDate = computed(() => {
  const d = auth.user?.created_at
  return d ? new Date(d).toLocaleDateString('en-US', { month: 'long', year: 'numeric' }) : '—'
})

const stats = computed(() => {
  const o = ordersStore.orders
  return [
    { label: 'Total Orders',  value: o.length,                                              color: 'text-gray-900' },
    { label: 'Delivered',     value: o.filter(x => x.status === 'delivered').length,        color: 'text-green-600' },
    { label: 'In Progress',   value: o.filter(x => ['pending','confirmed','shipped'].includes(x.status)).length, color: 'text-blue-600' },
    { label: 'Total Spent',   value: '$' + o.reduce((s, x) => s + x.total, 0).toFixed(2),  color: 'text-purple-600' }
  ]
})

function logout() {
  auth.logout()
  router.push('/login')
}

onMounted(async () => {
  await Promise.all([auth.fetchMe(), ordersStore.fetchOrders()])
})
</script>
