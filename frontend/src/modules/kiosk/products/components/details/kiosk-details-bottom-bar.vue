<!-- src/components/FixedBottomBar.vue -->
<template>
	<section class="bottom-14 left-0 fixed flex justify-center w-full">
		<div
			class="z-20 flex items-center gap-4 bg-white/50 shadow-2xl backdrop-blur-sm p-2 rounded-full"
		>
			<!-- Size Selection -->
			<div
				v-if="isMultipleSizes"
				class="flex items-center bg-gray-200 rounded-full overflow-x-auto no-scrollbar"
			>
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
				<Pencil
					v-if="displayIcon === 'update'"
					class="w-6 sm:w-8 h-6 sm:h-8"
				/>
				<Plus
					v-else
					class="w-6 sm:w-8 h-6 sm:h-8"
				/>
				<p class="text-xl sm:text-3xl">{{ formatPrice(totalPrice) }}</p>
			</button>
		</div>
	</section>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreProductSizeDTO } from '@/modules/admin/store-products/models/store-products.model'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import { Pencil, Plus } from 'lucide-vue-next'
import { computed } from "vue"

const {sizes, selectedSizeId, totalPrice, displayIcon = "add" } = defineProps<{
  sizes: StoreProductSizeDTO[]
  selectedSizeId: number | undefined
  totalPrice: number,
  displayIcon?: "add" | "update",
}>()

const emits = defineEmits<{
  (e: 'selectSize', size: StoreProductSizeDTO): void
  (e: 'addToCart'): void
}>()

const onSizeSelect = (size: StoreProductSizeDTO) => {
  emits('selectSize', size)
}

const handleAddToCart = () => {
  emits('addToCart')
}

const isMultipleSizes = computed(() => sizes.length > 1)

const isSelected = (sizeId: number) => {
  return selectedSizeId === sizeId
}
</script>

<style scoped></style>
