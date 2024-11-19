<template>
	<button
		@click="goBack"
		class="top-10 right-6 sm:right-10 z-10 fixed bg-gray-600 p-2 sm:p-3 rounded-full"
		aria-label="Go Back"
	>
		<X class="w-6 sm:w-7 h-6 sm:h-7 text-white" />
	</button>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.back()
}
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
