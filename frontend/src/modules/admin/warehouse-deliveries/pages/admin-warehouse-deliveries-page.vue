<template>
	<AdminWarehouseDeliveriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
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
import {Card, CardContent} from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import {usePaginationFilter} from '@/core/hooks/use-pagination-filter.hook'
import AdminWarehouseDeliveriesList
  from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-list.vue'
import AdminWarehouseDeliveriesToolbar
  from '@/modules/admin/warehouse-deliveries/components/list/admin-warehouse-deliveries-toolbar.vue'
import type {WarehouseDeliveryFilter} from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import {warehouseStocksService} from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import {useQuery} from '@tanstack/vue-query'
import {computed} from 'vue'
import {useHasRole} from "@/core/hooks/use-has-roles.hook";
import {EmployeeRole} from "@/modules/admin/employees/models/employees.models";

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<WarehouseDeliveryFilter>({})
const isRegion = useHasRole([EmployeeRole.REGION_WAREHOUSE_MANAGER])
const isWarehouse = useHasRole([EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER])

const { data: warehouseDeliveriesResponse, isLoading } = useQuery({
  queryKey: computed(() => ['warehouse-deliveries', filter.value]),
  queryFn: () => warehouseStocksService.getWarehouseDeliveries(filter.value),
  enabled: computed(() =>
    isWarehouse.value || (isRegion.value && Boolean(filter.value.warehouseId))
  )
})
</script>

<style scoped></style>
