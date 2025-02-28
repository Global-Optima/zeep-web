<template>
	<button
		@click="onUpdateSize(size)"
		:class="['rounded-full sm:w-24 sm:h-24 w-16 h-16', isSelected ? 'bg-primary text-white' : 'bg-gray-100 text-gray-800']"
	>
		<p class="font-medium text-xl sm:text-4xl">{{ size.name.charAt(0) }}</p>
		<p class="text-sm">{{ size.size }}</p>
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
