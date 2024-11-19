<template>
	<button
		@click="onUpdateSize(size)"
		:class="['rounded-full sm:w-20 sm:h-20 w-16 h-16', isSelected ? 'bg-primary text-primary-foreground' : '']"
	>
		<p class="text-xl sm:text-3xl">{{ size.name.charAt(0) }}</p>
	</button>
</template>

<script setup lang="ts">
import type { ProductSizeDTO } from '@/modules/kiosk/products/models/product.model'



const {size, isSelected} = defineProps<{size: ProductSizeDTO, isSelected: boolean}>()
const emit = defineEmits(["click:size"])

const onUpdateSize = (newSize: ProductSizeDTO) => {
  emit("click:size", newSize)
}
</script>

<style scoped></style>
