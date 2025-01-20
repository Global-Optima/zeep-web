<template>
	<AdminWarehouseStocksCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminWarehouseStocksCreateForm from '@/modules/admin/warehouse-stocks/components/create/admin-warehouse-stocks-create-form.vue'
import type { AddMultipleWarehouseStockDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: AddMultipleWarehouseStockDTO[]) => warehouseStocksService.addMultipleWarehouseStock(dto) ,
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-stocks'] })
		router.push({ name: getRouteName("ADMIN_WAREHOUSE_STOCKS") })
	},
})

function handleCreate(dto: AddMultipleWarehouseStockDTO[]) {
  if (dto.length === 0) return
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
