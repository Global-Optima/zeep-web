<template>
	<AdminStockMaterialCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminStockMaterialCategoriesList from '@/modules/admin/stock-material-categories/components/list/admin-stock-material-categories-list.vue'
import AdminStockMaterialCategoriesToolbar from '@/modules/admin/stock-material-categories/components/list/admin-stock-material-categories-toolbar.vue'
import type { StockMaterialCategoryFilterDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StockMaterialCategoryFilterDTO>({})

const { data: categoriesResponse } = useQuery({
  queryKey: computed(() => ['admin-stock-material-categories', filter.value]),
  queryFn: () => stockMaterialCategoryService.getAll(filter.value),
})

function updateFilter(updatedFilter: StockMaterialCategoryFilterDTO) {
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
