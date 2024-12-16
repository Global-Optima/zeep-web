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
		<CardFooter class="flex justify-center">
			<PaginationWithMeta
				v-if="storeStocksResponse"
				:meta="storeStocksResponse.pagination"
				@page-change="updatePage"
			/>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import AdminStoreStocksList from '@/modules/admin/store-warehouse/components/list/admin-store-stocks-list.vue'
import AdminStoreStocksToolbar from '@/modules/admin/store-warehouse/components/list/admin-store-stocks-toolbar.vue'
import type { StoreStocksFilter } from '@/modules/admin/store-warehouse/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-warehouse/services/store-stocks.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StoreStocksFilter>({
  page: 1,
  pageSize: 10,
})

const { currentStoreId } = useCurrentStoreStore()

// Derived reactive value for queryKey
const queryKey = computed(() => ['store-stocks', { storeId: currentStoreId, ...filter.value }])

const { data: storeStocksResponse } = useQuery({
  queryKey: queryKey,
  queryFn: () => {
    if (!currentStoreId) throw new Error('No store ID available')
    return storeStocksService.getStoreStocks(currentStoreId, filter.value)
  },
  enabled: computed(() => !!currentStoreId),
})

function updateFilter(updatedFilter: Partial<StoreStocksFilter>) {
  filter.value = updatedFilter
}

function updatePage(page: number) {
  updateFilter({ page })
}
</script>

<style scoped></style>
