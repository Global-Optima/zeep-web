<!-- src/components/FixedBottomBar.vue -->
<template>
	<section class="bottom-4 left-0 fixed flex justify-center w-full">
		<div
			class="z-20 flex items-center gap-4 bg-white/50 shadow-2xl backdrop-blur-sm px-6 py-4 rounded-full"
		>
			<!-- Size Selection -->
			<div class="flex items-center bg-gray-200 rounded-full overflow-x-auto no-scrollbar">
				<KioskDetailsSizes
					v-for="size in sizes"
					:key="size.id"
					:size="size"
					:is-selected="isSelected(size.id)"
					@click:size="onSizeSelect"
				/>
			</div>

			<!-- Add to Cart Button -->
			<button
				@click="handleAddToCart"
				class="flex items-center gap-3 bg-primary px-10 py-5 rounded-full text-primary-foreground"
			>
				<Plus class="w-6 sm:w-8 h-6 sm:h-8" />
				<p class="text-xl sm:text-3xl">{{ formattedPrice }}</p>
			</button>
		</div>
	</section>
</template>

<script setup lang="ts">
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import type { ProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { Plus } from 'lucide-vue-next'
import { computed } from 'vue'

const props = defineProps<{
  sizes: ProductSizeDTO[]
  selectedSizeId: number | undefined
  totalPrice: number
  formatPrice: (price: number) => string
}>()

const emits = defineEmits<{
  (e: 'selectSize', size: ProductSizeDTO): void
  (e: 'addToCart'): void
}>()

const onSizeSelect = (size: ProductSizeDTO) => {
  emits('selectSize', size)
}

const handleAddToCart = () => {
  emits('addToCart')
}

const isSelected = (sizeId: number) => {
  return props.selectedSizeId === sizeId
}

const formattedPrice = computed(() => props.formatPrice(props.totalPrice))
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
