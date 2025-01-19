<template>
	<div class="mx-auto max-w-6xl">
		<div class="gap-4 grid grid-cols-2 md:grid-cols-3">
			<div class="col-span-2">
				<AdminWarehouseStockRequestsDetailsMaterialsTable
					v-if="stockRequest"
					:stockRequest="stockRequest"
				/>
			</div>

			<div class="col-span-full md:col-span-1">
				<AdminWarehouseStockRequestsDetailsInfo
					v-if="stockRequest"
					:request="stockRequest"
				/>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import AdminWarehouseStockRequestsDetailsInfo from '@/modules/admin/warehouse-stock-requests/components/details/admin-warehouse-stock-requests-details-info.vue'
import AdminWarehouseStockRequestsDetailsMaterialsTable from '@/modules/admin/warehouse-stock-requests/components/details/admin-warehouse-stock-requests-details-materials-table.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const stockRequestId = route.params.id as string

const {
  data: stockRequest,
} = useQuery({
  queryKey: computed(() => ['stock-request', Number(stockRequestId)]),
  queryFn: () => stockRequestsService.getStockRequestById(Number(stockRequestId)),
  enabled: !isNaN(Number(stockRequestId)),
})
</script>
