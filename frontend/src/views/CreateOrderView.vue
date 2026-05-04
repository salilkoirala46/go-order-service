<template>
  <div class="max-w-lg mx-auto">
    <div class="flex items-center gap-3 mb-6">
      <button @click="router.back()" class="btn-secondary p-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <div>
        <h2 class="text-xl font-bold text-gray-900">New Order</h2>
        <p class="text-sm text-gray-500">Fill in the details below</p>
      </div>
    </div>

    <div class="card">
      <form @submit.prevent="submit" class="space-y-5">
        <div>
          <label class="label">Product name</label>
          <input v-model="form.product" type="text" class="input" placeholder="e.g. Widget Pro" required />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">Quantity</label>
            <input v-model.number="form.quantity" type="number" min="1" class="input" placeholder="1" required />
          </div>
          <div>
            <label class="label">Unit price ($)</label>
            <input v-model.number="form.price" type="number" min="0.01" step="0.01" class="input" placeholder="0.00" required />
          </div>
        </div>

        <!-- Order summary -->
        <div v-if="form.quantity > 0 && form.price > 0" class="bg-blue-50 rounded-lg p-4">
          <div class="flex justify-between text-sm text-gray-600 mb-1">
            <span>{{ form.quantity }} × ${{ form.price.toFixed(2) }}</span>
            <span>Subtotal</span>
          </div>
          <div class="flex justify-between font-bold text-gray-900">
            <span class="text-lg">${{ (form.quantity * form.price).toFixed(2) }}</span>
            <span>Total</span>
          </div>
        </div>

        <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

        <div class="flex gap-3 pt-2">
          <button type="button" @click="router.back()" class="btn-secondary flex-1">Cancel</button>
          <button type="submit" class="btn-primary flex-1" :disabled="loading">
            <svg v-if="loading" class="animate-spin -ml-1 mr-2 w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ loading ? 'Placing order…' : 'Place Order' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'

const store   = useOrdersStore()
const router  = useRouter()
const loading = ref(false)
const error   = ref('')
const form    = ref({ product: '', quantity: 1, price: 0 })

async function submit() {
  error.value   = ''
  loading.value = true
  try {
    const order = await store.createOrder(form.value.product, form.value.quantity, form.value.price)
    router.push(`/orders/${order.id}`)
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to create order. Try again.'
  } finally {
    loading.value = false
  }
}
</script>
