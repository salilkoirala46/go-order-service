<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 to-blue-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-14 h-14 bg-blue-500 rounded-2xl mb-4 shadow-lg">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-white">Create an account</h1>
        <p class="text-blue-200 text-sm mt-1">Start managing your orders today</p>
      </div>

      <div class="card shadow-xl">
        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="label">Full name</label>
            <input v-model="form.name" type="text" class="input" placeholder="Alice Smith" required />
          </div>
          <div>
            <label class="label">Email</label>
            <input v-model="form.email" type="email" class="input" placeholder="you@example.com" required />
          </div>
          <div>
            <label class="label">Password <span class="text-gray-400 font-normal">(min 6 chars)</span></label>
            <input v-model="form.password" type="password" class="input" placeholder="••••••••" minlength="6" required />
          </div>

          <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

          <button type="submit" class="btn-primary w-full" :disabled="loading">
            <svg v-if="loading" class="animate-spin -ml-1 mr-2 w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ loading ? 'Creating account…' : 'Create account' }}
          </button>
        </form>

        <p class="text-center text-sm text-gray-600 mt-4">
          Already have an account?
          <RouterLink to="/login" class="text-blue-600 hover:underline font-medium">Sign in</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth    = useAuthStore()
const router  = useRouter()
const loading = ref(false)
const error   = ref('')
const form    = ref({ name: '', email: '', password: '' })

async function submit() {
  error.value   = ''
  loading.value = true
  try {
    await auth.register(form.value.name, form.value.email, form.value.password)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Registration failed. Try again.'
  } finally {
    loading.value = false
  }
}
</script>
