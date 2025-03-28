<template>
	<!-- Toolbar -->
	<AdminStockRequestsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<!-- Main Content -->
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!stockRequestsResponse || stockRequestsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Запросы на склад не найдены
				</p>

				<AdminStockRequestsList
					v-else
					:requests="stockRequestsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="stockRequestsResponse"
					:meta="stockRequestsResponse.pagination"
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
import AdminStockRequestsList from '@/modules/admin/stock-requests/components/list/admin-stock-requests-list.vue'
import AdminStockRequestsToolbar from '@/modules/admin/stock-requests/components/list/admin-stock-requests-toolbar.vue'
import { type GetStockRequestsFilter } from '@/modules/admin/stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/stock-requests/services/stock-requests.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import {useHasRole} from "@/core/hooks/use-has-roles.hook";
import {EmployeeRole} from "@/modules/admin/employees/models/employees.models";

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<GetStockRequestsFilter>({})
const isFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])
const isRegion = useHasRole(EmployeeRole.REGION_WAREHOUSE_MANAGER)
const isWarehouse = useHasRole([EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER])
const isStore = useHasRole([EmployeeRole.BARISTA, EmployeeRole.STORE_MANAGER])

const { data: stockRequestsResponse, isLoading } = useQuery({
  queryKey: computed(() => ['stock-requests', filter.value]),
  queryFn: () => stockRequestsService.getStockRequests(filter.value),
  enabled: computed(() =>
    isStore.value ||
    isWarehouse.value ||
    (isFranchisee.value && Boolean(filter.value.storeId)) ||
    (isRegion.value && Boolean(filter.value.warehouseId))
  )
});
</script>
