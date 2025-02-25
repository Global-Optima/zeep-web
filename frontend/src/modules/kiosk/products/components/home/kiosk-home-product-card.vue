<template>
	<div
		class="flex flex-col justify-between bg-white p-6 rounded-[32px] h-full cursor-pointer"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<LazyImage
				src="https://www.nicepng.com/png/full/106-1060376_starbucks-iced-coffee-png-vector-library-pumpkin-spice.png"
				alt="Изображение товара"
				class="rounded-lg w-full h-44 sm:h-60 object-contain"
			/>

			<h3
				class="mt-6 text-base sm:text-xl line-clamp-2"
				data-testid="product-title"
			>
				{{ product.name }}
			</h3>
		</div>

		<div class="mt-4">
			<p
				class="font-medium text-primary text-xl sm:text-2xl"
				data-testid="product-price"
			>
				{{ formatPrice(product.storePrice) }}
			</p>
		</div>
	</div>
</template>

<script setup lang="ts">
 import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
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
