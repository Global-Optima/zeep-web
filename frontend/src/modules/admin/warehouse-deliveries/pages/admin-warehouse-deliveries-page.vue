<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminWarehouseDeliveriesList from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-list.vue'
import AdminWarehouseDeliveriesToolbar from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-toolbar.vue'
import type { WarehouseDeliveryFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<WarehouseDeliveryFilter>({})

const { data: warehouseDeliveriesResponse, isPending } = useQuery({
  queryKey: computed(() => ['warehouse-deliveries', filter.value]),
  queryFn: () => warehouseStocksService.getWarehouseDeliveries(filter.value),
})
</script>

<style scoped></style>
