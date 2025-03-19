<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminIngredientCategoryDetailsForm
		v-else
		:ingredientCategory="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminIngredientCategoryDetailsForm from '@/modules/admin/ingredient-categories/components/details/admin-ingredient-category-details-form.vue'
import type { UpdateIngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const categoryId = route.params.id as string

const { data: categoryDetails } = useQuery({
	queryKey: computed(() => ['admin-ingredient-categories-details', categoryId]),
	queryFn: () => ingredientsService.getIngredientCategoryById(Number(categoryId)),
	enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateIngredientCategoryDTO }) =>
		ingredientsService.updateIngredientCategory(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных категории ингредиентов. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories-details', categoryId] })
		toast({
			title: 'Успех!',
			description: 'Данные категории ингредиентов успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении категории ингредиентов.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateIngredientCategoryDTO) {
	if (isNaN(Number(categoryId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор категории.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(categoryId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
