<template>
	<div class="flex flex-col gap-10 px-8">
		<div
			v-for="category in visibleCategories"
			:key="category.id"
		>
			<!-- Category Name -->
			<p class="mb-6 font-medium text-3xl">{{ category.name }}</p>

			<!-- Additives Grid -->
			<div class="gap-4 grid grid-cols-1 md:grid-cols-2">
				<KioskDetailsAdditivesCard
					v-for="additive in category.visibleAdditives"
					:key="additive.additiveId"
					:additive="additive"
					:is-selected="isAdditiveSelected(category, additive.additiveId)"
					:has-category-price="category.hasAnyPrice"
					@click:additive="() => onAdditiveToggle(category, additive)"
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import type {
  StoreAdditiveCategoryDTO,
  StoreAdditiveCategoryItemDTO
} from '@/modules/admin/store-additives/models/store-additves.model'
import { computed } from 'vue'
import KioskDetailsAdditivesCard from './kiosk-details-additives-card.vue'

const props = defineProps<{
  categories: StoreAdditiveCategoryDTO[];
  selectedAdditives: Record<number, StoreAdditiveCategoryItemDTO[]>;
}>()

const emits = defineEmits<{
  (e: 'toggleAdditive', category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO): void
}>()

function onAdditiveToggle(
  category: StoreAdditiveCategoryDTO,
  additive: StoreAdditiveCategoryItemDTO
) {
  emits('toggleAdditive', category, additive)
}

function isAdditiveSelected(
  category: StoreAdditiveCategoryDTO,
  additiveId: number
): boolean {
  return props.selectedAdditives[category.id]?.some(
    (a) => a.additiveId === additiveId
  ) || false
}

const visibleCategories = computed(() => {
  return props.categories
    .map((category) => {
      const hasAnyPrice = category.additives.some((a) => a.storePrice > 0)

      const visibleAdditives = category.additives.filter((a) => !a.isHidden)

      return {
        ...category,
        hasAnyPrice,
        visibleAdditives,
      }
    })
    .filter((category) => category.visibleAdditives.length > 0)
})
</script>

<style scoped></style>
