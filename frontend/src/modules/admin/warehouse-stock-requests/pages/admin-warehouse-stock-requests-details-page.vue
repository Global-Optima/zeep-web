<template>
	<div class="mx-auto max-w-6xl">
		<div class="gap-4 grid grid-cols-2 md:grid-cols-3">
			<!-- Product Details Section -->
			<div class="col-span-2">
				<AdminWarehouseStockRequestsDetailsMaterialsTable
					v-if="stockRequest"
					:items="stockRequest.stockMaterials"
				/>
			</div>

			<!-- Media and Category Blocks Section -->
			<div class="col-span-full md:col-span-1">
				<AdminWarehouseStockRequestsDetailsInfo
					v-if="stockRequest"
					:request="stockRequest"
					@update:status="onUpdateStatus"
				/>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast'
import type { StockRequestStatus } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import AdminWarehouseStockRequestsDetailsInfo from '@/modules/admin/warehouse-stock-requests/components/details/admin-warehouse-stock-requests-details-info.vue'
import AdminWarehouseStockRequestsDetailsMaterialsTable from '@/modules/admin/warehouse-stock-requests/components/details/admin-warehouse-stock-requests-details-materials-table.vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const queryClient = useQueryClient()
const {toast} = useToast()

const storeStockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
  queryKey: computed(() => ['warehouse-stock-request', storeStockRequestId]),
	queryFn: () => stockRequestsService.getStockRequestById(Number(storeStockRequestId)),
  enabled: !isNaN(Number(storeStockRequestId)),
})

// const {mutate: updateStatusMutation} = useMutation({
// 		mutationFn: (data: {id: number, dto: UpdateStockRequestStatusDTO}) => stockRequestService.updateStockRequestStatus(data.id, data.dto),
// 		onSuccess: () => {
//       toast({title: "Статус успешно обновлен"})
//       queryClient.invalidateQueries({ queryKey: ['warehouse-stock-requests'] })
//       queryClient.invalidateQueries({ queryKey: ['warehouse-stock-request', storeStockRequestId] })

// 		},
// 		onError: () => {
// 			toast({title: "Произошла ошибка при обновлении"})
// 		},
// })

const onUpdateStatus = (status: StockRequestStatus) => {
}
</script>
