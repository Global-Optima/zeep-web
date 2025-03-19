<template>
	<AdminWarehouseDeliveriesCreateForm
		@on-submit="handleCreate"
		@on-close="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminWarehouseDeliveriesCreateForm from '@/modules/admin/warehouse-deliveries/components/create/admin-warehouse-deliveries-create-form.vue'
import type { ReceiveWarehouseDelivery } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: ReceiveWarehouseDelivery) =>
		warehouseStocksService.receiveWarehouseDelivery(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Добавление новой поставки на склад. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-deliveries'] })
		queryClient.invalidateQueries({ queryKey: ['warehouse-stocks'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Поставка успешно добавлена на склад.',
		})
		router.push({ name: getRouteName('ADMIN_WAREHOUSE_DELIVERIES') })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении поставки на склад.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: ReceiveWarehouseDelivery) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
