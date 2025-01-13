<template>
	<AdminWarehouseStocksToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!warehouseStocksResponse || warehouseStocksResponse.data.levels.length === 0"
				class="text-muted-foreground"
			>
				Складские запасы не найдены
			</p>

			<AdminWarehouseStocksList
				v-else
				:stocks="warehouseStocksResponse.data.levels"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="warehouseStocksResponse"
				:meta="warehouseStocksResponse.pagination"
				@update:page="updatePage"
				@update:pageSize="updatePageSize"
			/>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminWarehouseStocksList from '@/modules/admin/warehouse-stocks/components/list/admin-warehouse-stocks-list.vue'
import AdminWarehouseStocksToolbar from '@/modules/admin/warehouse-stocks/components/list/admin-warehouse-stocks-toolbar.vue'
import type { GetWarehouseStockFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<GetWarehouseStockFilter>({
})

const { data: warehouseStocksResponse } = useQuery({
  queryKey: computed(() => ['warehouse-stocks', filter.value]),
  queryFn: () => warehouseStocksService.getWarehouseStocks(filter.value),
})

function updateFilter(updatedFilter: GetWarehouseStockFilter) {
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
