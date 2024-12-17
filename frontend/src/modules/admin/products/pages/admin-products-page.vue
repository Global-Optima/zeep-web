<template>
	<AdminProductsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<AdminProductsList :products="products" />
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminProductsList from '@/modules/admin/products/components/list/admin-products-list.vue'
import AdminProductsToolbar from '@/modules/admin/products/components/list/admin-products-toolbar.vue'
import type { SuppliersFilter } from '@/modules/admin/suppliers/models/suppliers.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import type { StoresFilter } from '@/modules/stores/models/stores-dto.model'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<SuppliersFilter>({
  searchTerm: '',
})

const { data: products } = useQuery({
  queryKey: computed(() => ['suppliers', filter.value]),
  queryFn: () => productsService.getProducts(filter.value),
  initialData: []
})

function updateFilter(updatedFilter: Partial<StoresFilter>) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>

<style scoped></style>
