<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
		<AdminProductCategoriesToolbar
			:filter="filter"
			@update:filter="updateFilter"
		/>

		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!productCategoriesResponse || productCategoriesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Категории товаров не найдены
				</p>
				<AdminProductCategoriesList
					v-else
					:productCategories="productCategoriesResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="productCategoriesResponse"
					:meta="productCategoriesResponse.pagination"
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
import AdminProductCategoriesList from '@/modules/admin/product-categories/components/list/admin-product-categories-list.vue'
import AdminProductCategoriesToolbar from '@/modules/admin/product-categories/components/list/admin-product-categories-toolbar.vue'
import type { ProductCategoriesFilterDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<ProductCategoriesFilterDTO>({});

const { data: productCategoriesResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-product-categories', filter.value]),
  queryFn: () => productsService.getAllProductCategories(filter.value),
})
</script>

<style scoped></style>
