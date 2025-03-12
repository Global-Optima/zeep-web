<!-- src/components/AdditivesSection.vue -->
<template>
	<div :class="cn('flex flex-col gap-8 px-8', containerClass)">
		<template
			v-for="category in categories"
			:key="category.id"
		>
			<div v-if="category.additives.length > 0">
				<p class="font-medium text-3xl">{{ category.name }}</p>
				<div class="gap-4 grid grid-cols-2 mt-5">
					<KioskDetailsAdditivesCard
						v-for="additive in category.additives"
						:key="additive.additiveId"
						:additive="additive"
						:is-default="additive.isDefault"
						:is-selected="isAdditiveSelected(category, additive.additiveId)"
						@click:additive="() => onAdditiveToggle(category, additive)"
					/>
				</div>
			</div>
		</template>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import KioskDetailsAdditivesCard from '@/modules/kiosk/products/components/details/kiosk-details-additives-card.vue'

defineProps<{
  categories: StoreAdditiveCategoryDTO[]
  isAdditiveSelected: (category: StoreAdditiveCategoryDTO, additiveId: number) => boolean
  containerClass?: string
}>()

const emits = defineEmits<{
  (e: 'toggleAdditive', category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO): void
}>()

const onAdditiveToggle = (category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO) => {
  emits('toggleAdditive', category, additive)
}
</script>

<style scoped></style>
