<template>
	<div class="bg-white rounded-3xl p-6 flex flex-col justify-between h-full">
		<router-link :to="{ name: getRouteName('KIOSK_PRODUCT_DETAILS'), params: { id: product.id } }">
			<div>
				<img
					:src="product.image"
					alt="Product Image"
					class="w-full h-44 object-contain rounded-lg"
				/>

				<h3 class="text-base sm:text-xl mt-3 line-clamp-2 min-h-[3rem]">{{ product.title }}</h3>
			</div>

			<div class="mt-4">
				<p class="text-xl sm:text-2xl font-medium">{{ formatPrice(product.price) }}</p>
			</div>
		</router-link>
	</div>
</template>

<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { Products } from '@/modules/products/models/product.model'

const {product} = defineProps<{product:Products}>()
</script>

<style scoped></style>
