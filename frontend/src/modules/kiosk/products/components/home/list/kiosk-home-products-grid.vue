<template>
	<div>
		<div v-if="isLoading">
			<KioskHomeProductsLoading :skeletonCount="skeletonCount" />
		</div>
		<div
			v-else-if="!isLoading && products.length === 0"
			class="flex justify-center items-center h-20 text-gray-500"
		>
			<p class="text-3xl">{{ emptyMessage }}</p>
		</div>
		<div
			v-else
			class="gap-4 grid grid-cols-2 sm:grid-cols-3"
		>
			<KioskHomeProductCard
				v-for="product in products"
				:key="product.id"
				:product="product"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeProductsLoading from '@/modules/kiosk/products/components/home/list/kiosk-home-products-loading.vue'
import { defineProps } from 'vue'

const props = defineProps<{
  products: StoreProductDTO[]
  isLoading: boolean
  emptyMessage?: string
  skeletonCount?: number
}>()

const emptyMessage = props.emptyMessage || 'Ничего не найдено'
const skeletonCount = props.skeletonCount || 6
</script>
