<!-- src/components/AdditivesSection.vue -->
<template>
	<div :class="cn('flex flex-col gap-6 mt-2 px-8', containerClass)">
		<template
			v-for="category in categories"
			:key="category.id"
		>
			<div v-if="category.additives.length > 0">
				<p class="font-medium text-lg sm:text-2xl">{{ category.name }}</p>
				<div class="flex flex-wrap gap-2 mt-4">
					<KioskDetailsAdditivesCard
						v-for="additive in category.additives"
						:key="additive.id"
						:additive="additive"
						:is-default="additive.isDefault"
						:is-selected="isAdditiveSelected(category.id, additive.id)"
						@click:additive="() => onAdditiveToggle(category.id, additive)"
					/>
				</div>
			</div>
		</template>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { AdditiveCategoryItemDTO } from '@/modules/admin/additives/models/additives.model'
import type { StoreAdditiveCategoryDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import KioskDetailsAdditivesCard from '@/modules/kiosk/products/components/details/kiosk-details-additives-card.vue'

defineProps<{
  categories: StoreAdditiveCategoryDTO[]
  isAdditiveSelected: (categoryId: number, additiveId: number) => boolean
  containerClass?: string
}>()

const emits = defineEmits<{
  (e: 'toggleAdditive', categoryId: number, additive: AdditiveCategoryItemDTO): void
}>()

const onAdditiveToggle = (categoryId: number, additive: AdditiveCategoryItemDTO) => {
  emits('toggleAdditive', categoryId, additive)
}
</script>

<style scoped></style>
