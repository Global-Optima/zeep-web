<template>
	<AdminStoreStocksDetailsForm
		v-if="storeStockData"
		:initialData="storeStockData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreStocksDetailsForm from '@/modules/admin/store-warehouse/components/details/admin-store-stocks-details-form.vue'
import type { UpdateStoreStock } from '@/modules/admin/store-warehouse/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-warehouse/services/store-stocks.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const storeStockId = route.params.id as string

const queryClient = useQueryClient()

const {currentStoreId} = useCurrentStoreStore()


const { data: storeStockData } = useQuery({
  queryKey: computed(() => ['store-stock', storeStockId, { storeId: currentStoreId}]),
	queryFn: () =>{    if (!currentStoreId) throw new Error('No store ID available')
  return storeStocksService.getStoreStock(currentStoreId, Number(storeStockId))},
  enabled: computed(() => !!currentStoreId),
})

const updateMutation = useMutation({
	mutationFn: (updatedData: UpdateStoreStock) => {
    if (!currentStoreId) throw new Error('No store ID available')
    return storeStocksService.updateStoreStock(currentStoreId, Number(storeStockId), updatedData)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-stocks'] })
		queryClient.invalidateQueries({ queryKey: ['store-stock', storeStockId] })
		router.push({ name: getRouteName('ADMIN_STORE_STOCKS') })
	},
})

function handleUpdate(updatedData: UpdateStoreStock) {
	updateMutation.mutate(updatedData)
}

function handleCancel() {
	router.back()
}
</script>
