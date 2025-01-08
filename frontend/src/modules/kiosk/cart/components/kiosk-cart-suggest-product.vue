<template>
	<div class="flex flex-col justify-between bg-white p-4 sm:p-6 rounded-2xl h-full">
		<div>
			<img
				:src="product.imageUrl"
				alt="Product Image"
				class="rounded-lg w-full h-44 object-contain"
			/>

			<h3 class="mt-3 text-sm sm:text-xl">{{ product.name }}</h3>
		</div>

		<div class="flex justify-between items-center mt-3">
			<p class="text-lg sm:text-2xl">{{ 1200 }}</p>
			<button
				class="bg-gray-200 p-2 sm:p-3 rounded-full text-white"
				aria-label="Добавить к заказу"
			>
				<!-- TODO: ADD SUGGEST FOOD
      @click="() => cartStore.addToCart(product, {
            id: 4,
            name: 'S',
            basePrice: 1200,
            measure: 'гр.'
        }, [])" -->

				<Icon
					icon="mingcute:add-line"
					class="text-black text-xl"
				/>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import type { StoreProductDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { Icon } from '@iconify/vue'

const { product } = defineProps<{ product: StoreProductDetailsDTO }>()

const cartStore = useCartStore()
</script>

<style scoped></style>
