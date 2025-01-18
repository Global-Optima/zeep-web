<template>
	<div class="mx-auto max-w-6xl">
		<div class="gap-4 grid grid-cols-2 md:grid-cols-3">
			<div class="col-span-2">
				<AdminStoreStockRequestsDetailsMaterialsTable
					v-if="stockRequest"
					:stockRequest="stockRequest"
				/>
			</div>

			<div class="col-span-full md:col-span-1">
				<AdminStoreStockRequestsDetailsInfo
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
import AdminStoreStockRequestsDetailsInfo from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-info.vue'
import AdminStoreStockRequestsDetailsMaterialsTable from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-materials-table.vue'
import type { StockRequestStatus } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const queryClient = useQueryClient()
const {toast} = useToast()

const storeStockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
  queryKey: computed(() => ['stock-request', storeStockRequestId]),
	queryFn: () => stockRequestsService.getStockRequestById(Number(storeStockRequestId)),
  enabled: !isNaN(Number(storeStockRequestId)),
})

// const {mutate: updateStatusMutation} = useMutation({
// 		mutationFn: (data: {id: number, dto: UpdateStoreStockRequestStatusDTO}) => stockRequestsService.updateStockRequestIngredients(data.id, data.dto),
// 		onSuccess: () => {
//       toast({title: "Статус успешно обновлен"})
//       queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
//       queryClient.invalidateQueries({ queryKey: ['stock-request', storeStockRequestId] })

// 		},
// 		onError: () => {
// 			toast({title: "Произошла ошибка при обновлении"})
// 		},
// })

const onUpdateStatus = (status: StockRequestStatus) => {
  // updateStatusMutation({id: Number(storeStockRequestId), dto: {status }})
}
</script>
