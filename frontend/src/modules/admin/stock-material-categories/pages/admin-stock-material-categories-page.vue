<template>
	<AdminStockMaterialCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!categoriesResponse || categoriesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Категории не найдены
				</p>
				<AdminStockMaterialCategoriesList
					v-else
					:categories="categoriesResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="categoriesResponse"
					:meta="categoriesResponse.pagination"
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
import AdminStockMaterialCategoriesList from '@/modules/admin/stock-material-categories/components/list/admin-stock-material-categories-list.vue'
import AdminStockMaterialCategoriesToolbar from '@/modules/admin/stock-material-categories/components/list/admin-stock-material-categories-toolbar.vue'
import type { StockMaterialCategoryFilterDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<StockMaterialCategoryFilterDTO>({});

const { data: categoriesResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-stock-material-categories', filter.value]),
  queryFn: () => stockMaterialCategoryService.getAll(filter.value),
})
</script>

<style scoped></style>
