<template>
	<KioskHomeProductsGrid
		:products="products"
		:isLoading="isLoading"
		emptyMessage="Ничего не найдено"
	/>
</template>

<script setup lang="ts">
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import KioskHomeProductsGrid from '@/modules/kiosk/products/components/home/list/kiosk-home-products-grid.vue'
import { defineProps } from 'vue'

defineProps<{
  products: StoreProductDTO[]
  isLoading: boolean
}>()
</script>
