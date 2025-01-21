<template>
	<AdminWarehouseDeliveriesCreateForm
		@on-submit="handleCreate"
		@on-close="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminWarehouseDeliveriesCreateForm from '@/modules/admin/warehouse-deliveries/components/create/admin-warehouse-deliveries-create-form.vue'
import type { ReceiveWarehouseDelivery } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: ReceiveWarehouseDelivery) => warehouseStocksService.receiveWarehouseDelivery(dto) ,
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-deliveries'] })
		router.push({ name: getRouteName("ADMIN_WAREHOUSE_DELIVERIES") })
	},
})

function handleCreate(dto: ReceiveWarehouseDelivery){
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
