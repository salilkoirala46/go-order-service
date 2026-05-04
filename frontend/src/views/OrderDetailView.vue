<template>
  <div class="max-w-2xl mx-auto space-y-5">
    <!-- Back -->
    <div class="flex items-center gap-3">
      <button @click="router.back()" class="btn-secondary p-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <div>
        <h2 class="text-xl font-bold text-gray-900">Order #{{ route.params.id }}</h2>
        <p class="text-sm text-gray-500">Order details and status management</p>
      </div>
    </div>

    <LoadingSpinner v-if="store.loading" />

    <template v-else-if="store.current">
      <!-- Summary card -->
      <div class="card">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h3 class="text-lg font-bold text-gray-900">{{ store.current.product }}</h3>
            <p class="text-sm text-gray-400 mt-0.5">Placed on {{ formatDate(store.current.created_at) }}</p>
          </div>
          <StatusBadge :status="store.current.status" />
        </div>

        <dl class="grid grid-cols-2 gap-4 mt-4 border-t border-gray-100 pt-4">
          <div>
            <dt class="text-xs text-gray-500 uppercase tracking-wide">Quantity</dt>
            <dd class="mt-1 font-semibold text-gray-900">{{ store.current.quantity }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-500 uppercase tracking-wide">Unit Price</dt>
            <dd class="mt-1 font-semibold text-gray-900">${{ store.current.price.toFixed(2) }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-500 uppercase tracking-wide">Total</dt>
            <dd class="mt-1 text-xl font-bold text-blue-600">${{ store.current.total.toFixed(2) }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-500 uppercase tracking-wide">Last Updated</dt>
            <dd class="mt-1 font-semibold text-gray-900">{{ formatDate(store.current.updated_at) }}</dd>
          </div>
        </dl>
      </div>

      <!-- Status timeline -->
      <div class="card">
        <h4 class="font-semibold text-gray-900 mb-4">Order Progress</h4>
        <ol class="relative border-l-2 border-gray-200 ml-3 space-y-4">
          <li
            v-for="step in statusSteps"
            :key="step.value"
            class="ml-5"
          >
            <span
              :class="[
                'absolute -left-2.5 w-5 h-5 rounded-full border-2 flex items-center justify-center',
                isPast(step.value) ? 'bg-blue-600 border-blue-600' :
                  isCurrent(step.value) ? 'bg-white border-blue-600' : 'bg-white border-gray-300'
              ]"
            >
              <svg v-if="isPast(step.value)" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </span>
            <p :class="['font-medium text-sm', isPast(step.value) || isCurrent(step.value) ? 'text-gray-900' : 'text-gray-400']">
              {{ step.label }}
            </p>
            <p class="text-xs text-gray-400">{{ step.description }}</p>
          </li>
        </ol>
      </div>

      <!-- Update status -->
      <div v-if="store.current.status !== 'delivered' && store.current.status !== 'cancelled'" class="card">
        <h4 class="font-semibold text-gray-900 mb-3">Update Status</h4>
        <div class="flex gap-2 flex-wrap">
          <button
            v-for="next in nextStatuses"
            :key="next.value"
            @click="updateStatus(next.value)"
            :disabled="updating"
            :class="['btn text-sm', next.danger ? 'btn-danger' : 'btn-primary']"
          >
            {{ updating ? 'Updating…' : next.label }}
          </button>
        </div>
        <p v-if="updateError" class="text-sm text-red-600 mt-2">{{ updateError }}</p>
      </div>
    </template>

    <div v-else class="card text-center py-10 text-gray-500">
      Order not found.
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import StatusBadge from '@/components/StatusBadge.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const store       = useOrdersStore()
const route       = useRoute()
const router      = useRouter()
const updating    = ref(false)
const updateError = ref('')

const statusOrder = ['pending', 'confirmed', 'shipped', 'delivered']

const statusSteps = [
  { value: 'pending',   label: 'Order Placed',  description: 'Awaiting confirmation' },
  { value: 'confirmed', label: 'Confirmed',      description: 'Order is being prepared' },
  { value: 'shipped',   label: 'Shipped',        description: 'On the way to you' },
  { value: 'delivered', label: 'Delivered',      description: 'Order completed' }
]

const transitions = {
  pending:   [{ value: 'confirmed', label: 'Confirm Order' }, { value: 'cancelled', label: 'Cancel', danger: true }],
  confirmed: [{ value: 'shipped',   label: 'Mark Shipped'  }, { value: 'cancelled', label: 'Cancel', danger: true }],
  shipped:   [{ value: 'delivered', label: 'Mark Delivered' }]
}

const nextStatuses = computed(() => transitions[store.current?.status] ?? [])

function isPast(val) {
  const cur = statusOrder.indexOf(store.current?.status)
  const idx = statusOrder.indexOf(val)
  return idx < cur || store.current?.status === 'delivered'
}

function isCurrent(val) {
  return store.current?.status === val
}

async function updateStatus(status) {
  updateError.value = ''
  updating.value    = true
  try {
    await store.updateStatus(Number(route.params.id), status)
  } catch (e) {
    updateError.value = e.response?.data?.error ?? 'Update failed'
  } finally {
    updating.value = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => store.fetchOrder(Number(route.params.id)))
</script>
