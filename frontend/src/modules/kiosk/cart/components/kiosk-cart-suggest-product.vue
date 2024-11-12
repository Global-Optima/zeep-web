<template>
	<div class="bg-white rounded-2xl p-4 sm:p-6 flex flex-col justify-between h-full">
		<div>
			<img
				:src="product.imageUrl"
				alt="Product Image"
				class="w-full h-44 object-contain rounded-lg"
			/>

			<h3 class="text-sm sm:text-xl mt-3">{{ product.name }}</h3>
		</div>

		<div class="flex justify-between items-center mt-3">
			<p class="text-lg sm:text-2xl">{{ formatPrice(product.basePrice) }}</p>
			<button
				class="bg-gray-200 text-white p-2 sm:p-3 rounded-full"
				aria-label="Добавить к заказу"
				@click="() => cartStore.addToCart(product, {
            id: 4,
            name: 'S',
            basePrice: product.basePrice,
            measure: 'гр.'
        }, [])"
			>
				<Icon
					icon="mingcute:add-line"
					class="text-xl text-black"
				/>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import type { StoreProductDetails } from '@/modules/kiosk/products/models/product.model'
import { Icon } from '@iconify/vue'

const { product } = defineProps<{ product: StoreProductDetails }>()

const cartStore = useCartStore()
</script>

<style scoped></style>
