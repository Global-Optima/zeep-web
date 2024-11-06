<template>
	<div
		class="bg-white rounded-3xl p-6 flex flex-col justify-between h-full cursor-pointer"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<img
				:src="product.image"
				alt="Product Image"
				class="w-full h-44 object-contain rounded-lg"
				data-testid="product-image"
			/>

			<h3
				class="text-base sm:text-xl mt-3 line-clamp-2 min-h-[3rem]"
				data-testid="product-title"
			>
				{{ product.title }}
			</h3>
		</div>

		<div class="mt-4">
			<p
				class="text-xl sm:text-2xl font-medium"
				data-testid="product-price"
			>
				{{ formatPrice(product.startPrice) }}
			</p>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { defineProps, defineEmits } from 'vue';
 import { formatPrice } from '@/core/utils/price.utils';
 import type { Products } from '@/modules/products/models/product.model';

 const props = defineProps<{
product: Products;
 }>();

 const emit = defineEmits<{
(e: 'select-product', product: Products): void;
 }>();

 const selectProduct = () => {
emit('select-product', props.product);
 };
</script>

<style scoped></style>
