<template>
	<div>
		<!-- Toolbar -->
		<AdminWarehouseStockRequestsToolbar
			:filter="filter"
			@update:filter="updateFilter"
		/>

		<!-- Main Content -->
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!stockRequestsResponse || stockRequestsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Запросы на склад не найдены
				</p>

				<AdminWarehouseStockRequestsList
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
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import { WAREHOUSE_STOCK_REQUEST_STATUSES, type GetStoreStockRequestsFilter } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { storeStockRequestService } from '@/modules/admin/store-stock-requests/services/store-stock-request.service'
import AdminWarehouseStockRequestsList from '@/modules/admin/warehouse-stock-requests/components/list/admin-warehouse-stock-requests-list.vue'
import AdminWarehouseStockRequestsToolbar from '@/modules/admin/warehouse-stock-requests/components/list/admin-warehouse-stock-requests-toolbar.vue'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const {currentEmployee} = useEmployeeAuthStore()
const filter = ref<GetStoreStockRequestsFilter>({
  warehouseId: currentEmployee?.warehouseDetails?.warehouseId,
  statuses: WAREHOUSE_STOCK_REQUEST_STATUSES
});

const { data: stockRequestsResponse } = useQuery({
  queryKey: computed(() => ['warehouse-stock-requests', filter.value]),
  queryFn: () => storeStockRequestService.getStockRequests(filter.value),
  enabled: computed(() => Boolean(currentEmployee?.warehouseDetails?.warehouseId))
});

function updateFilter(updatedFilter: GetStoreStockRequestsFilter) {
  filter.value = { ...filter.value, ...updatedFilter };
}

function updatePage(page: number) {
  updateFilter({ page, pageSize: DEFAULT_PAGINATION_META.pageSize });
}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize, page: DEFAULT_PAGINATION_META.page });
}
</script>
