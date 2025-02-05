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
				{{ formatPrice(product.storePrice) }}
			</p>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { formatPrice } from '@/core/utils/price.utils'
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { useCurrentProductStore } from '@/modules/kiosk/products/components/hooks/use-current-product.hook'

// const router = useRouter()
const currentProductStore = useCurrentProductStore()

 const {product} = defineProps<{
  product: StoreProductDTO;
}>();

 const selectProduct = () => {
  currentProductStore.openModal(product.id)
};
</script>

<style scoped></style>
