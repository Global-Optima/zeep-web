<template>
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminProductCategoriesList from '@/modules/admin/product-categories/components/list/admin-product-categories-list.vue'
import AdminProductCategoriesToolbar from '@/modules/admin/product-categories/components/list/admin-product-categories-toolbar.vue'
import type { ProductCategoriesFilterDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<ProductCategoriesFilterDTO>({})

const { data: productCategoriesResponse } = useQuery({
  queryKey: computed(() => ['admin-product-categories', filter.value]),
  queryFn: () => productsService.getAllProductCategories(filter.value),
})

function updateFilter(updatedFilter: ProductCategoriesFilterDTO) {
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
