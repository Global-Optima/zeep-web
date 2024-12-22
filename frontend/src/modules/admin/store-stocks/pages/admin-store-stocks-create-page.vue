<template>
	<AdminStoreStocksCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreStocksCreateForm from '@/modules/admin/store-stocks/components/create/admin-store-stocks-create-form.vue'
import type { CreateMultipleStoreStock } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const {currentStoreId} = useCurrentStoreStore()

const createMutation = useMutation({
	mutationFn: (dto: CreateMultipleStoreStock) => {
    if (!currentStoreId) throw new Error('No store ID available')
    return storeStocksService.createMultipleStoreStock(currentStoreId, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-stocks'] })
		router.push({ name: getRouteName("ADMIN_STORE_STOCKS") })
	},
})

function handleCreate(dto: CreateMultipleStoreStock) {
  if (dto.ingredientStocks.length === 0) return
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
