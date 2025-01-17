<template>
	<aside
		class="sm:flex flex-col hidden bg-white px-6 py-6 rounded-3xl max-w-[250px] h-full overflow-y-auto no-scrollbar"
	>
		<ul class="flex flex-col justify-center">
			<li
				v-for="category in categories"
				:key="category.id"
				@click="onUpdateCategory(category.id)"
				:class="[
          'cursor-pointer text-2xl py-4 category-item',
          category.id === selectedCategoryId ? 'text-primary' : 'text-gray-600',
        ]"
			>
				{{ category.name }}
			</li>
		</ul>
	</aside>
</template>

<script setup lang="ts">
import type { ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'


const {selectedCategoryId, categories}= defineProps<{selectedCategoryId: number | null, categories: ProductCategoryDTO[]}>()

const emit= defineEmits({
  "update:category": null
})

const onUpdateCategory = (categoryId: number ) => {
  emit("update:category", categoryId)
}
</script>

<style scoped>
.category-item:not(:last-child) {
  border-bottom: 1px solid #f1f1f1;
}
</style>
