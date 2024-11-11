<template>
	<button
		@click="onUpdateSize(size)"
		:class="['rounded-full sm:w-24 sm:h-24 w-16 h-16', selectedSize === size ? 'bg-primary text-primary-foreground' : 'bg-white text-black']"
	>
		<p class="text-lg sm:text-2xl font-medium">{{ size.label }}</p>
		<!-- <p class="text-xs sm:text-base">{{ size.volume }} ml</p> -->
	</button>
</template>

<script setup lang="ts">
import type { ProductSizes } from '@/modules/products/models/product.model'

  const {size, selectedSize} = defineProps<{size: ProductSizes, selectedSize: ProductSizes}>()
  const emit = defineEmits(["click:size"])

  const onUpdateSize = (newSize: ProductSizes) => {
    emit("click:size", newSize)
  }
</script>

<style scoped></style>
