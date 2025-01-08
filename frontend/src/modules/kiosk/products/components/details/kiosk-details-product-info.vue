<!-- src/components/ProductInfo.vue -->
<template>
	<section :class="cn('p-8', containerClass)">
		<h1 class="font-medium text-2xl sm:text-4xl">{{ name }}</h1>
		<p class="mt-3 text-base sm:text-2xl">{{ description }}</p>

		<div class="mt-8">
			<KioskDetailsIngredients :ingredients='ingredients' :energy='energy' />
		</div>
	</section>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import KioskDetailsIngredients from '@/modules/kiosk/products/components/details/kiosk-details-ingredients.vue'

export interface  KioskDetailsEnergyDTO {
  ccal: number
  proteins: number
  carbs: number
  fats: number
}

defineProps<{
  name: string
  description: string
  containerClass?: string
  energy: KioskDetailsEnergyDTO
  ingredients: IngredientsDTO[]
}>()
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
