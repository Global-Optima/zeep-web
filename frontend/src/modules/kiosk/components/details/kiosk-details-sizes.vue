<template>
	<button
		@click="onUpdateSize(size)"
		:class="['rounded-full bg-white text-black sm:w-28 sm:h-28 w-[72px] h-[72px]', selectedSize === size ? 'border-2 border-primary' : '']"
	>
		<p class="text-base sm:text-2xl font-medium">{{ size.label }}</p>
		<p class="text-xs sm:text-base">{{ size.volume }} ml</p>
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
