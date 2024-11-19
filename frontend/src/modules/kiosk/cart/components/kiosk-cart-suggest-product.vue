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
			<p class="text-lg sm:text-2xl">{{ 1200 }}</p>
			<button
				class="bg-gray-200 text-white p-2 sm:p-3 rounded-full"
				aria-label="Добавить к заказу"
				@click="() => cartStore.addToCart(product, {
            id: 4,
            name: 'S',
            basePrice: 1200,
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
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import type { StoreProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { Icon } from '@iconify/vue'

const { product } = defineProps<{ product: StoreProductDetailsDTO }>()

const cartStore = useCartStore()
</script>

<style scoped></style>
