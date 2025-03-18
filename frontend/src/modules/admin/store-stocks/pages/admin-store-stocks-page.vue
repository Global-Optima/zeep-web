<template>
	<AdminStoreStocksToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storeStocksResponse || storeStocksResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Складские запасы не найдены
				</p>

				<AdminStoreStocksList
					v-else
					:stocks="storeStocksResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="storeStocksResponse"
					:meta="storeStocksResponse.pagination"
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
import AdminStoreStocksList from '@/modules/admin/store-stocks/components/list/admin-store-stocks-list.vue'
import AdminStoreStocksToolbar from '@/modules/admin/store-stocks/components/list/admin-store-stocks-toolbar.vue'
import type { GetStoreWarehouseStockFilterQuery } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<GetStoreWarehouseStockFilterQuery>({})

const { data: storeStocksResponse, isLoading } = useQuery({
  queryKey: computed(() => ['store-stocks', filter.value]),
  queryFn: () => storeStocksService.getStoreWarehouseStockList(filter.value),
  enabled: computed(() => Boolean(filter.value.storeId))
})
</script>

<style scoped></style>
