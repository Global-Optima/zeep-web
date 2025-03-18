<template>
	<div
		class="flex flex-col justify-between bg-white p-6 rounded-[32px] h-full cursor-pointer"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<LazyImage
				:src="product.imageUrl"
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
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { useRouter } from 'vue-router'

// const router = useRouter()
// const currentProductStore = useCurrentProductStore()

const router = useRouter()

 const {product} = defineProps<{
  product: StoreProductDTO;
}>();

 const selectProduct = () => {
  // currentProductStore.openModal(product.id)
  router.push({name: getRouteName("KIOSK_PRODUCT"), params: {id: product.id}})
};
</script>

<style scoped></style>
