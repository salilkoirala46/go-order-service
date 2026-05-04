<template>
  <span :class="['inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium', colorClass]">
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ status: { type: String, required: true } })

const map = {
  pending:   { label: 'Pending',   color: 'bg-yellow-100 text-yellow-800' },
  confirmed: { label: 'Confirmed', color: 'bg-blue-100 text-blue-800'     },
  shipped:   { label: 'Shipped',   color: 'bg-purple-100 text-purple-800' },
  delivered: { label: 'Delivered', color: 'bg-green-100 text-green-800'   },
  cancelled: { label: 'Cancelled', color: 'bg-red-100 text-red-800'       }
}

const entry     = computed(() => map[props.status] ?? { label: props.status, color: 'bg-gray-100 text-gray-800' })
const label     = computed(() => entry.value.label)
const colorClass = computed(() => entry.value.color)
</script>
