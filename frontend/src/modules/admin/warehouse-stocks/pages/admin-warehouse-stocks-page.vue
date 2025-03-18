<template>
	<AdminWarehouseStocksToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!warehouseStocksResponse || warehouseStocksResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Складские запасы не найдены
				</p>

				<AdminWarehouseStocksList
					v-else
					:stocks="warehouseStocksResponse.data"
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminWarehouseStocksList from '@/modules/admin/warehouse-stocks/components/list/admin-warehouse-stocks-list.vue'
import AdminWarehouseStocksToolbar from '@/modules/admin/warehouse-stocks/components/list/admin-warehouse-stocks-toolbar.vue'
import type { GetWarehouseStockFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<GetWarehouseStockFilter>({})

const { data: warehouseStocksResponse, isLoading } = useQuery({
  queryKey: computed(() => ['warehouse-stocks', filter.value]),
  queryFn: () => warehouseStocksService.getWarehouseStocks(filter.value),
  enabled: computed(() => Boolean(filter.value.warehouseId))
})
</script>

<style scoped></style>
