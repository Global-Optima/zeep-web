<template>
	<div class="mb-16">
		<h2
			:ref="el => setCategoryRef(category.id, el)"
			class="my-8 px-4 text-4xl"
		>
			{{ category.name }}
		</h2>
		<KioskHomeProductsGrid
			:products="products"
			:isLoading="isLoading"
		/>
	</div>
</template>

<script setup lang="ts">
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import KioskHomeProductsGrid from '@/modules/kiosk/products/components/home/list/kiosk-home-products-grid.vue'
import type { ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { type ComponentPublicInstance } from 'vue'

defineProps<{
  category: ProductCategoryDTO
  products: StoreProductDTO[]
  isLoading: boolean
  setCategoryRef: (id: number, el: Element | ComponentPublicInstance | null) => void
}>()
</script>
