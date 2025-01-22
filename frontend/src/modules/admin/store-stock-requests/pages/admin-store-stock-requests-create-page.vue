<template>
	<AdminStoreStockRequestsCreateForm
		@submit="handleCreate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreStockRequestsCreateForm from '@/modules/admin/store-stock-requests/components/create/admin-store-stock-requests-create-form.vue'
import type { CreateStockRequestDTO, StockRequestStockMaterialDTO } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const {toast} = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockRequestDTO) => {
    return stockRequestsService.createStockRequest(dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
		router.push({ name: getRouteName("ADMIN_STORE_STOCK_REQUESTS") })
	},
  onError:(error: AxiosError<{error: string}>) => {
      const message =  error.response?.data.error ?? "Ошибка при создании" // TODO: reconsider the error messages (localization)
      toast({description: message, variant: "destructive"})
  }
})

function handleCreate(items: StockRequestStockMaterialDTO[]) {
  const dto: CreateStockRequestDTO = {
    stockMaterials: items
  }

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
