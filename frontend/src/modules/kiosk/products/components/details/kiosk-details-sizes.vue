<template>
	<button
		@click="onUpdateSize(size)"
		:class="['rounded-full md:size-28 size-16', isSelected ? 'bg-primary text-white' : 'bg-gray-100 text-gray-800']"
	>
		<p class="font-medium text-xl sm:text-4xl">{{ size.name.charAt(0) }}</p>
		<p class="mt-1 text-base">{{ size.size }} {{ size.unit.name.toLowerCase() }}</p>
	</button>
</template>

<script setup lang="ts">
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'

const {size, isSelected} = defineProps<{size: StoreProductSizeDetailsDTO, isSelected: boolean}>()
const emit = defineEmits(["click:size"])

const onUpdateSize = (newSize: StoreProductSizeDetailsDTO) => {
  emit("click:size", newSize)
}
</script>

<style scoped></style>
