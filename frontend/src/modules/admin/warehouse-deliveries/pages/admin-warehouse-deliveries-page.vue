<template>
	<AdminWarehouseDeliveriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!warehouseDeliveriesResponse || warehouseDeliveriesResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Складские запасы не найдены
			</p>

			<AdminWarehouseDeliveriesList
				v-else
				:deliveries="warehouseDeliveriesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="warehouseDeliveriesResponse"
				:meta="warehouseDeliveriesResponse.pagination"
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
import AdminWarehouseDeliveriesList from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-list.vue'
import AdminWarehouseDeliveriesToolbar from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-toolbar.vue'
import type { WarehouseDeliveryFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<WarehouseDeliveryFilter>({
})

const { data: warehouseDeliveriesResponse } = useQuery({
  queryKey: computed(() => ['warehouse-deliveries', filter.value]),
  queryFn: () => warehouseStocksService.getWarehouseDeliveries(filter.value),
})

function updateFilter(updatedFilter: WarehouseDeliveryFilter) {
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
