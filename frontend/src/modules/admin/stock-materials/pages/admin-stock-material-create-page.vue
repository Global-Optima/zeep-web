<template>
	<AdminStockMaterialCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStockMaterialCreateForm from '@/modules/admin/stock-materials/components/create/admin-stock-material-create-form.vue'
import type { CreateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockMaterialDTO) => stockMaterialsService.createStockMaterial(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового материала. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-materials'] })
		toast({
			title: 'Успех!',
			description: 'Материал успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании материала.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateStockMaterialDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
