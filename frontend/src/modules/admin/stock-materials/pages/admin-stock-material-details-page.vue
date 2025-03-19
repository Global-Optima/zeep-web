<template>
	<p v-if="!additiveDetails">Товар не найден</p>

	<AdminStockMaterialDetailsForm
		v-else
		:stockMaterial="additiveDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
		:readonly="!canUpdate"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminStockMaterialDetailsForm from '@/modules/admin/stock-materials/components/details/admin-stock-material-details-form.vue'
import type { UpdateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const stockMaterialId = route.params.id as string

const canUpdate = useHasRole([EmployeeRole.ADMIN])

const { data: additiveDetails } = useQuery({
	queryKey: ['admin-stock-material-details', stockMaterialId],
	queryFn: () => stockMaterialsService.getStockMaterialById(Number(stockMaterialId)),
	enabled: !isNaN(Number(stockMaterialId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStockMaterialDTO }) =>
		stockMaterialsService.updateStockMaterial(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных товара. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-materials'] })
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-details', stockMaterialId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные товара успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных товара.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateStockMaterialDTO) {
	if (isNaN(Number(stockMaterialId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор товара.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(stockMaterialId), dto: data })
}

function handleCancel() {
	router.back()
}
</script>
