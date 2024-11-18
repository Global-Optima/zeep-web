<!-- src/components/ProductInfo.vue -->
<template>
	<section :class="cn('p-8', containerClass)">
		<h1 class="font-medium text-2xl sm:text-4xl">{{ name }}</h1>
		<p class="mt-3 text-base sm:text-xl">{{ description }}</p>

		<div class="mt-8">
			<KioskDetailsEnergy :energy="energy" />
		</div>
	</section>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import KioskDetailsEnergy from '@/modules/kiosk/products/components/details/kiosk-details-energy.vue'

defineProps<{
  name: string
  description: string
  energy: {
    ccal: number
    proteins: number
    carbs: number
    fats: number
  }
  containerClass?: string
}>()
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
