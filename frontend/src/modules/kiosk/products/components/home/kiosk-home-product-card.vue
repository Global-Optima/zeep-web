<template>
	<div
		class="relative flex flex-col justify-between bg-white p-8 rounded-[52px] h-full transition-all duration-300"
		:class="{
      'cursor-pointer': !product.isOutOfStock,
      'cursor-not-allowed opacity-60 !border-primary' : product.isOutOfStock
    }"
		@click="selectProduct"
		data-testid="product-card"
	>
		<div>
			<LazyImage
				:src="product.imageUrl"
				alt="Изображение товара"
				class="rounded-[32px] w-full h-64 sm:h-80 object-contain"
			/>

			<h3
				class="mt-6 text-base sm:text-2xl line-clamp-2"
				data-testid="product-title"
			>
				{{ product.name }}
			</h3>
		</div>

		<div class="flex flex-wrap justify-between items-start gap-4 mt-10">
			<p
				class="font-medium text-xl sm:text-3xl"
				:class="product.isOutOfStock ? 'text-gray-400' : 'text-primary'"
				data-testid="product-price"
			>
				{{ formatPrice(product.storePrice) }}
			</p>

			<p
				v-if="product.isOutOfStock"
				class="top-8 right-8 absolute bg-slate-200 px-4 py-2 rounded-3xl text-slate-800 text-base"
			>
				Нет в наличии
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
   if (product.isOutOfStock) {
     return
   }
  router.push({name: getRouteName("KIOSK_PRODUCT"), params: {id: product.id}})
};
</script>

<style scoped></style>
