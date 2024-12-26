<template>
	<AdminStoreStocksToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminStoreStocksList from '@/modules/admin/store-stocks/components/list/admin-store-stocks-list.vue'
import AdminStoreStocksToolbar from '@/modules/admin/store-stocks/components/list/admin-store-stocks-toolbar.vue'
import type { StoreStocksFilter } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StoreStocksFilter>({
  page: 1,
  pageSize: 10,
})

const { currentStoreId } = useCurrentStoreStore()

const { data: storeStocksResponse } = useQuery({
  queryKey: computed(() => ['store-stocks', { storeId: currentStoreId, ...filter.value }]),
  queryFn: () => {
    if (!currentStoreId) throw new Error('No store ID available')
    return storeStocksService.getStoreStocks(currentStoreId, filter.value)
  },
  enabled: computed(() => !!currentStoreId),
})

function updateFilter(updatedFilter: StoreStocksFilter) {
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
