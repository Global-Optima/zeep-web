<template>
	<AdminWarehouseStocksDetailsForm
		v-if="stockData"
		:initialData="stockData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminWarehouseStocksDetailsForm from '@/modules/admin/warehouse-stocks/components/details/admin-warehouse-stocks-details-form.vue'
import type { UpdateWarehouseStockDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const warehouseStockId = route.params.id as string

const queryClient = useQueryClient()

const { data: stockData } = useQuery({
  queryKey: computed(() => ['warehouse-stock', warehouseStockId]),
	queryFn: () => warehouseStocksService.getWarehouseStockById(Number(warehouseStockId)),
  enabled: !isNaN(Number(warehouseStockId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}: {id:number, dto: UpdateWarehouseStockDTO}) => warehouseStocksService.updateWarehouseStocksById(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-stocks'] })
		queryClient.invalidateQueries({ queryKey: ['warehouse-stock', warehouseStockId] })
		router.push({ name: getRouteName('ADMIN_WAREHOUSE_STOCKS') })
	},
})

function handleUpdate(updatedData: UpdateWarehouseStockDTO) {
  if (isNaN(Number(warehouseStockId))) return router.back()
	updateMutation.mutate({id: Number(warehouseStockId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
