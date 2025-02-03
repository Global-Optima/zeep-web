<template>
	<AdminStockRequestsCreateForm
		@submit="handleCreate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminStockRequestsCreateForm from '@/modules/admin/stock-requests/components/create/admin-stock-requests-create-form.vue'
import type { CreateStockRequestDTO, StockRequestStockMaterialDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/stock-requests/services/stock-requests.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockRequestDTO) => stockRequestsService.createStockRequest(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Запрос на склад создается. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
		toast({
			title: 'Успех!',
			description: 'Запрос на склад успешно создан.',
		})
		router.push({ name: getRouteName('ADMIN_STORE_STOCK_REQUESTS') })
	},
	onError: (error: AxiosError<{ error: string }>) => {
		const message = error.response?.data.error ?? 'Ошибка при создании запроса на склад.'
		toast({
			title: 'Ошибка',
			description: message,
			variant: 'destructive',
		})
	},
})

function handleCreate(items: StockRequestStockMaterialDTO[]) {
	if (items.length === 0) {
		toast({
			title: 'Ошибка',
			description: 'Список материалов пуст. Пожалуйста, добавьте материалы перед отправкой.',
			variant: 'destructive',
		})
		return
	}

	const dto: CreateStockRequestDTO = {
		stockMaterials: items,
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
