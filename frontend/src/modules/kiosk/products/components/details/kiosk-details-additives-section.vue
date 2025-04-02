<!-- src/components/AdditivesSection.vue -->
<template>
	<div :class="cn('flex flex-col gap-10 px-8', containerClass)">
		<div
			v-for="category in visibleCategories"
			:key="category.id"
		>
			<p class="mb-6 font-medium text-3xl">{{ category.name }}</p>

			<div class="gap-4 grid grid-cols-1 md:grid-cols-2">
				<KioskDetailsAdditivesCard
					v-for="additive in category.additives"
					:key="additive.additiveId"
					:additive="additive"
					:is-selected="isAdditiveSelected(category, additive.additiveId)"
					@click:additive="() => onAdditiveToggle(category, additive)"
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import KioskDetailsAdditivesCard from '@/modules/kiosk/products/components/details/kiosk-details-additives-card.vue'
import { computed } from "vue"

const props = defineProps<{
  categories: StoreAdditiveCategoryDTO[]
  isAdditiveSelected: (category: StoreAdditiveCategoryDTO, additiveId: number) => boolean
  containerClass?: string
}>()

const emits = defineEmits<{
  (e: 'toggleAdditive', category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO): void
}>()

const onAdditiveToggle = (
  category: StoreAdditiveCategoryDTO,
  additive: StoreAdditiveCategoryItemDTO
) => {
  emits('toggleAdditive', category, additive)
}

const visibleCategories = computed(() => {
  return props.categories
    .map(category => {
      const visibleAdditives = category.additives.filter(a => !a.isHidden)
      return {
        ...category,
        additives: visibleAdditives,
      }
    })
    .filter(category => category.additives.length > 0)
})
</script>

<style scoped></style>
