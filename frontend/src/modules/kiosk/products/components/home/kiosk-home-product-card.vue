<template>
	<div
		class="bg-white rounded-3xl p-6 flex flex-col justify-between h-full cursor-pointer"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<img
				:src="product.imageUrl"
				alt="Product Image"
				class="w-full h-44 object-contain rounded-lg"
				data-testid="product-image"
			/>

			<h3
				class="text-base sm:text-xl mt-3 line-clamp-2 min-h-[3rem]"
				data-testid="product-title"
			>
				{{ product.name }}
			</h3>
		</div>

		<div class="mt-4">
			<p
				class="text-xl sm:text-2xl font-medium"
				data-testid="product-price"
			>
				{{ formatPrice(product.basePrice) }}
			</p>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { formatPrice } from '@/core/utils/price.utils'
import type { StoreProducts } from '@/modules/kiosk/products/models/product.model'

 const props = defineProps<{
product: StoreProducts;
 }>();

 const emit = defineEmits<{
  (e: 'select-product', product: number): void;
 }>();

 const selectProduct = () => {
  emit('select-product', props.product.id);
 };
</script>

<style scoped></style>
