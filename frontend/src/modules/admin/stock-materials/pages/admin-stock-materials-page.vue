<template>
	<AdminStockMaterialsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!stockMaterialsResponse || stockMaterialsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Складские товары не найдены
			</p>
			<AdminStockMaterialsList
				v-else
				:stock-materials="stockMaterialsResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="stockMaterialsResponse"
				:meta="stockMaterialsResponse.pagination"
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
import AdminStockMaterialsList from '@/modules/admin/stock-materials/components/list/admin-stock-materials-list.vue'
import AdminStockMaterialsToolbar from '@/modules/admin/stock-materials/components/list/admin-stock-materials-toolbar.vue'
import type { StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StockMaterialsFilter>({})

const { data: stockMaterialsResponse } = useQuery({
  queryKey: computed(() => ['admin-stock-materials', filter.value]),
  queryFn: () => stockMaterialsService.getAllStockMaterials(filter.value),
})

function updateFilter(updatedFilter: StockMaterialsFilter) {
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
