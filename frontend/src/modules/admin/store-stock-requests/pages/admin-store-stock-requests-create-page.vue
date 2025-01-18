<template>
	<AdminStoreStockRequestsCreateForm
		@submit="handleCreate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreStockRequestsCreateForm from '@/modules/admin/store-stock-requests/components/create/admin-store-stock-requests-create-form.vue'
import type { CreateStoreStockRequestDTO, CreateStoreStockRequestItemDTO } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { storeStockRequestService } from '@/modules/admin/store-stock-requests/services/store-stock-request.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()


const createMutation = useMutation({
	mutationFn: (dto: CreateStoreStockRequestDTO) => {
    return storeStockRequestService.createStockRequest(dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
		router.push({ name: getRouteName("ADMIN_STORE_STOCK_REQUESTS") })
	},
})

function handleCreate(items: CreateStoreStockRequestItemDTO[]) {
  const dto: CreateStoreStockRequestDTO = {
    items: items
  }

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
