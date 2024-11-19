<!-- src/components/AdditivesSection.vue -->
<template>
	<div :class="cn('flex flex-col gap-6 mt-2 px-8', containerClass)">
		<div
			v-for="category in categories"
			:key="category.id"
		>
			<p class="font-medium text-lg sm:text-2xl">{{ category.name }}</p>
			<div class="flex flex-wrap gap-2 mt-4">
				<KioskDetailsAdditivesCard
					v-for="additive in category.additives"
					:key="additive.id"
					:additive="additive"
					:is-default="isAdditiveDefault(additive.id)"
					:is-selected="isAdditiveSelected(category.id, additive.id)"
					@click:additive="() => onAdditiveToggle(category.id, additive)"
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import KioskDetailsAdditivesCard from '@/modules/kiosk/products/components/details/kiosk-details-additives-card.vue'
import type { AdditiveCategoryDTO, AdditiveDTO } from '@/modules/kiosk/products/models/product.model'

defineProps<{
  categories: AdditiveCategoryDTO[]
  isAdditiveDefault: (additiveId: number) => boolean
  isAdditiveSelected: (categoryId: number, additiveId: number) => boolean
  containerClass?: string
}>()

const emits = defineEmits<{
  (e: 'toggleAdditive', categoryId: number, additive: AdditiveDTO): void
}>()

const onAdditiveToggle = (categoryId: number, additive: AdditiveDTO) => {
  emits('toggleAdditive', categoryId, additive)
}
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
