<template>
	<aside
		class="hidden sm:flex py-8 px-6 bg-white h-full overflow-hidden flex-col max-w-[250px] rounded-3xl"
	>
		<p class="text-3xl font-semibold">ZEEP</p>
		<ul class="flex flex-col justify-center gap-4 mt-8">
			<li
				v-for="category in categories"
				:key="category.id"
				@click="onUpdateCategory(category.id)"
				:class="[
						'cursor-pointer text-xl',
						category.id === selectedCategoryId ? 'text-primary' : 'text-gray-600',
					]"
			>
				{{ category.name }}
			</li>
		</ul>
	</aside>
</template>

<script setup lang="ts">
import type { ProductCategory } from '@/modules/kiosk/products/models/product.model'

const {selectedCategoryId, categories}= defineProps<{selectedCategoryId: number | null, categories: ProductCategory[]}>()

const emit= defineEmits({
  "update:category": null
})

const onUpdateCategory = (categoryId: number ) => {
  emit("update:category", categoryId)
}
</script>

<style scoped></style>
