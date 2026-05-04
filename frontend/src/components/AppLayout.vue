<template>
  <div class="flex h-screen bg-gray-50">
    <AppSidebar />
    <div class="flex flex-col flex-1 min-w-0">
      <AppNavbar />
      <main class="flex-1 overflow-y-auto p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
import AppSidebar from './AppSidebar.vue'
import AppNavbar from './AppNavbar.vue'
</script>
