<template>
	<AdminIngredientCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminIngredientCategoryCreateForm from '@/modules/admin/ingredient-categories/components/create/admin-ingredient-category-create-form.vue'
import type { CreateIngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateIngredientCategoryDTO) => ingredientsService.createIngredientCategory(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой категории ингредиентов. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Категория ингредиентов успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании категории ингредиентов.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateIngredientCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
