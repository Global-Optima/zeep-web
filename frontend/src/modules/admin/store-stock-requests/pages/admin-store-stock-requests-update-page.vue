<template>
	<AdminStoreStockRequestsUpdateForm
		v-if="stockRequest"
		:initialData="stockRequest.stockMaterials"
		@submit="handleUpdate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreStockRequestsUpdateForm, { type StockRequestItemForm } from '@/modules/admin/store-stock-requests/components/update/admin-store-stock-requests-update-form.vue'
import type { StockRequestStockMaterialDTO } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const route = useRoute()
const stockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
  queryKey: computed(() => ['stock-request', Number(stockRequestId)]),
	queryFn: () => stockRequestsService.getStockRequestById(Number(stockRequestId)),
  enabled: !isNaN(Number(stockRequestId)),
})

const updateMutation = useMutation({
	mutationFn: (props: {id: number, dto: StockRequestStockMaterialDTO[]}) => stockRequestsService.updateStockRequestMaterials(props.id, props.dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
    queryClient.invalidateQueries({ queryKey: ['stock-request', Number(stockRequestId)] })
		router.back()
	},
})

function handleUpdate(dto: StockRequestItemForm[]) {
	updateMutation.mutate({id: Number(stockRequestId), dto: dto.map(item => ({
    stockMaterialId: item.stockMaterialId,
	  quantity: item.quantity,
  }))})
}

function handleCancel() {
	router.back()
}
</script>
