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
import { toastError, toastSuccess } from '@/core/config/toast.config'
import AdminStoreStockRequestsDetailsInfo from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-info.vue'
import AdminStoreStockRequestsDetailsMaterialsTable from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-materials-table.vue'
import type { StoreStockRequestStatus, UpdateStoreStockRequestStatusDTO } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { storeStockRequestService } from '@/modules/admin/store-stock-requests/services/store-stock-request.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const queryClient = useQueryClient()

const storeStockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
  queryKey: computed(() => ['stock-request', storeStockRequestId]),
	queryFn: () => storeStockRequestService.getStockRequestById(Number(storeStockRequestId)),
  enabled: !isNaN(Number(storeStockRequestId)),
})

const {mutate: updateStatusMutation} = useMutation({
		mutationFn: (data: {id: number, dto: UpdateStoreStockRequestStatusDTO}) => storeStockRequestService.updateStockRequestStatus(data.id, data.dto),
		onSuccess: () => {
      toastSuccess("Статус успешно обновлен")
      queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
      queryClient.invalidateQueries({ queryKey: ['stock-request', storeStockRequestId] })

		},
		onError: () => {
			toastError("Произошла ошибка при обновлении")
		},
})

const onUpdateStatus = (status: StoreStockRequestStatus) => {
  updateStatusMutation({id: Number(storeStockRequestId), dto: {status }})
}
</script>
