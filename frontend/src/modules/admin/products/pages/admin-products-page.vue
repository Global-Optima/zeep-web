<template>
	<AdminProductsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!productsResponse || productsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Продукты не найдены
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminProductsList from '@/modules/admin/products/components/list/admin-products-list.vue'
import AdminProductsToolbar from '@/modules/admin/products/components/list/admin-products-toolbar.vue'
import type { ProductsFilterDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use the pagination filter hook
const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<ProductsFilterDTO>({});

const { data: productsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-products', filter.value]),
  queryFn: () => productsService.getProducts(filter.value),
})
</script>

<style scoped></style>
