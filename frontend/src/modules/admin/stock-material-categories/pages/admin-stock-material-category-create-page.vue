<template>
	<AdminStockMaterialCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStockMaterialCategoryCreateForm from '@/modules/admin/stock-material-categories/components/create/admin-stock-material-category-create-form.vue'
import type { CreateStockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockMaterialCategoryDTO) => stockMaterialCategoryService.create(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой категории материалов. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Категория материалов успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании категории материалов.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateStockMaterialCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
