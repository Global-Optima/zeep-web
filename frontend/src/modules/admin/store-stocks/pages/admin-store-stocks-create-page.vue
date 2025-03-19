<template>
	<AdminStoreStocksCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreStocksCreateForm from '@/modules/admin/store-stocks/components/create/admin-store-stocks-create-form.vue'
import type { AddMultipleStoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: AddMultipleStoreWarehouseStockDTO) =>
		storeStocksService.addMultipleStoreWarehouseStock(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Добавление новых запасов кафе. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-stocks'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Запасы кафе успешно добавлены.',
		})

    router.push({ name: getRouteName('ADMIN_STORE_STOCKS') })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении запасов кафе.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: AddMultipleStoreWarehouseStockDTO) {
	if (dto.ingredientStocks.length === 0) {
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
