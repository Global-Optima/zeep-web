<template>
	<AdminStoreStocksDetailsForm
		v-if="storeStockData"
		:initialData="storeStockData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreStocksDetailsForm from '@/modules/admin/store-stocks/components/details/admin-store-stocks-details-form.vue'
import type { UpdateStoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import { storeStocksService } from '@/modules/admin/store-stocks/services/store-stocks.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const storeStockId = route.params.id as string

const { data: storeStockData } = useQuery({
	queryKey: computed(() => ['store-stock', storeStockId]),
	queryFn: () => storeStocksService.getStoreWarehouseStockById(Number(storeStockId)),
	enabled: !isNaN(Number(storeStockId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStoreWarehouseStockDTO }) =>
		storeStocksService.updateStoreWarehouseStockById(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных запасов магазина. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-stocks'] })
		queryClient.invalidateQueries({ queryKey: ['store-stock', storeStockId] })
		toast({
			title: 'Успех!',
			description: 'Данные запасов магазина успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных запасов.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStoreWarehouseStockDTO) {
	if (isNaN(Number(storeStockId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор запаса.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(storeStockId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
