<template>
	<div
		class="flex flex-col justify-between bg-white p-6 rounded-3xl h-full cursor-pointer"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<img
				:src="product.imageUrl"
				alt="Product Image"
				class="rounded-lg w-full h-44 sm:h-60 object-contain"
				data-testid="product-image"
			/>

			<h3
				class="mt-3 line-clamp-2 min-h-[3rem] text-base sm:text-xl"
				data-testid="product-title"
			>
				{{ product.name }}
			</h3>
		</div>

		<div class="mt-4">
			<p
				class="font-medium text-xl sm:text-2xl"
				data-testid="product-price"
			>
				{{ formatPrice(product.basePrice) }}
			</p>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreProducts } from '@/modules/kiosk/products/models/product.model'
import { useRouter } from 'vue-router'

const router = useRouter()

 const {product} = defineProps<{
  product: StoreProducts;
}>();

 const selectProduct = () => {
  // productStore.selectProduct(product.id)
  router.push({name: getRouteName("KIOSK_DETAILS"), params: {id: product.id}})
};
</script>

<style scoped></style>
