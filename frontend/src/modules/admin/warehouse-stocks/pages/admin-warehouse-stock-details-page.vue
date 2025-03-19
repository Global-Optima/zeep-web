<template>
	<AdminWarehouseStocksDetailsForm
		v-if="stockData"
		:initialData="stockData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
		:readonly="!canUpdate"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminWarehouseStocksDetailsForm from '@/modules/admin/warehouse-stocks/components/details/admin-warehouse-stocks-details-form.vue'
import type { UpdateWarehouseStockDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const warehouseStockId = route.params.id as string

const canUpdate = useHasRole([EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER])

const { data: stockData } = useQuery({
	queryKey: computed(() => ['warehouse-stock', warehouseStockId]),
	queryFn: () => warehouseStocksService.getWarehouseStockById(Number(warehouseStockId)),
	enabled: !isNaN(Number(warehouseStockId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateWarehouseStockDTO }) =>
		warehouseStocksService.updateWarehouseStocksById(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных запаса склада. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-stocks'] })
		queryClient.invalidateQueries({ queryKey: ['warehouse-stock', warehouseStockId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные запаса склада успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных запаса склада.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateWarehouseStockDTO) {
	if (isNaN(Number(warehouseStockId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор запаса.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(warehouseStockId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
