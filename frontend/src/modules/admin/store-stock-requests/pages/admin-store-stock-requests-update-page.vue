<template>
	<AdminStoreStockRequestsUpdateForm
		v-if="stockRequest"
		:initialData="stockRequest.items"
		@submit="handleUpdate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreStockRequestsUpdateForm, { type StoreStockRequestItemForm } from '@/modules/admin/store-stock-requests/components/update/admin-store-stock-requests-update-form.vue'
import type { CreateStoreStockRequestItemDTO } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { storeStockRequestService } from '@/modules/admin/store-stock-requests/services/store-stock-request.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const route = useRoute()
const stockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
  queryKey: computed(() => ['stock-request', stockRequestId]),
	queryFn: () => storeStockRequestService.getStockRequestById(Number(stockRequestId)),
  enabled: !isNaN(Number(stockRequestId)),
})

const updateMutation = useMutation({
	mutationFn: (props: {id: number, dto: CreateStoreStockRequestItemDTO[]}) => storeStockRequestService.updateStockRequestIngredients(props.id, props.dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
    queryClient.invalidateQueries({ queryKey: ['stock-request', stockRequestId] })
		router.back()
	},
})

function handleUpdate(dto: StoreStockRequestItemForm[]) {
	updateMutation.mutate({id: Number(stockRequestId), dto: dto.map(item => ({
    stockMaterialId: item.stockMaterialId,
	  quantity: item.quantity,
  }))})
}

function handleCancel() {
	router.back()
}
</script>
