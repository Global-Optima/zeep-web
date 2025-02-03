<template>
	<AdminStockRequestsUpdateForm
		v-if="stockRequest"
		:initialData="stockRequest.stockMaterials"
		@submit="handleUpdate"
		@cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStockRequestsUpdateForm from '@/modules/admin/stock-requests/components/update/admin-stock-requests-update-form.vue'
import type { StockRequestStockMaterialDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/stock-requests/services/stock-requests.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const stockRequestId = route.params.id as string

const { data: stockRequest } = useQuery({
	queryKey: computed(() => ['stock-request', Number(stockRequestId)]),
	queryFn: () => stockRequestsService.getStockRequestById(Number(stockRequestId)),
	enabled: !isNaN(Number(stockRequestId)),
})

const updateMutation = useMutation({
	mutationFn: (props: { id: number; dto: StockRequestStockMaterialDTO[] }) =>
		stockRequestsService.updateStockRequestMaterials(props.id, props.dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление материалов запроса на склад. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
		queryClient.invalidateQueries({ queryKey: ['stock-request', Number(stockRequestId)] })
		toast({
			title: 'Успех!',
			description: 'Материалы запроса на склад успешно обновлены.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении материалов запроса на склад.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(dto: StockRequestStockMaterialDTO[]) {
	if (!dto.length) {
		toast({
			title: 'Ошибка',
			description: 'Список материалов пуст. Пожалуйста, добавьте материалы перед сохранением.',
			variant: 'destructive',
		})
		return
	}

	updateMutation.mutate({
		id: Number(stockRequestId),
		dto
	})
}

function handleCancel() {
	router.back()
}
</script>
