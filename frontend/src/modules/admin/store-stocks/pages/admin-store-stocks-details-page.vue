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
import AdminStoreStocksDetailsForm from '@/modules/admin/store-stocks/components/details/admin-store-stocks-details-form.vue'
import type { UpdateStoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const storeStockId = route.params.id as string

const queryClient = useQueryClient()

const { data: storeStockData } = useQuery({
  queryKey: computed(() => ['store-stock', storeStockId]),
	queryFn: () =>storeStocksService.getStoreWarehouseStockById(Number(storeStockId)),
  enabled: !isNaN(Number(storeStockId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}: {id:number, dto: UpdateStoreWarehouseStockDTO}) => storeStocksService.updateStoreWarehouseStockById(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-stocks'] })
		queryClient.invalidateQueries({ queryKey: ['store-stock', storeStockId] })
		router.push({ name: getRouteName('ADMIN_STORE_STOCKS') })
	},
})

function handleUpdate(updatedData: UpdateStoreWarehouseStockDTO) {
  if (isNaN(Number(storeStockId))) return router.back()
	updateMutation.mutate({id: Number(storeStockId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
