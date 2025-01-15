<template>
	<AdminProductsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!productsResponse || productsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Товары не найдены
			</p>
			<AdminProductsList
				v-else
				:products="productsResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="productsResponse"
				:meta="productsResponse.pagination"
				@update:page="updatePage"
				@update:pageSize="updatePageSize"
			/>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminProductsList from '@/modules/admin/products/components/list/admin-products-list.vue'
import AdminProductsToolbar from '@/modules/admin/products/components/list/admin-products-toolbar.vue'
import type { ProductsFilterDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<ProductsFilterDTO>({})

const { data: productsResponse } = useQuery({
  queryKey: computed(() => ['admin-products', filter.value]),
  queryFn: () => productsService.getProducts(filter.value),
})

function updateFilter(updatedFilter: ProductsFilterDTO) {
  filter.value = {...filter.value, ...updatedFilter}
}

function updatePage(page: number) {
  updateFilter({ pageSize: DEFAULT_PAGINATION_META.pageSize, page: page})

}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize: pageSize, page: DEFAULT_PAGINATION_META.page})
}
</script>

<style scoped></style>
