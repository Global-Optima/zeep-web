<template>
	<AdminIngredientCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminIngredientCreateForm from '@/modules/admin/ingredients/components/create/admin-ingredient-create-form.vue'
import type { CreateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateIngredientDTO) => ingredientsService.createIngredient(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового ингредиента. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredients'] })
		toast({
			title: 'Успех!',
			description: 'Ингредиент успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании ингредиента.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateIngredientDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
