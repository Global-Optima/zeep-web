<template>
	<AdminWarehouseStocksCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminWarehouseStocksCreateForm from '@/modules/admin/warehouse-stocks/components/create/admin-warehouse-stocks-create-form.vue'
import type { AddMultipleWarehouseStockDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: AddMultipleWarehouseStockDTO[]) => warehouseStocksService.addMultipleWarehouseStock(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Добавление новых запасов склада. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-stocks'] })
		toast({
			title: 'Успех!',
			description: 'Запасы склада успешно добавлены.',
		})
		router.push({ name: getRouteName('ADMIN_WAREHOUSE_STOCKS') })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении запасов склада.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: AddMultipleWarehouseStockDTO[]) {
	if (dto.length === 0) {
		toast({
			title: 'Ошибка',
			description: 'Список запасов пуст. Пожалуйста, добавьте запасы перед сохранением.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
