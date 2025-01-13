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
import { getRouteName } from '@/core/config/routes.config'
import AdminIngredientCategoryDetailsForm from '@/modules/admin/ingredient-categories/components/details/admin-ingredient-category-details-form.vue'
import type { UpdateIngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const categoryId = route.params.id as string

const queryClient = useQueryClient()

const { data: categoryDetails } = useQuery({
  queryKey: computed(() => ['admin-ingredient-categories-details', categoryId]),
	queryFn: () => ingredientsService.getIngredientCategoryById(Number(categoryId)),
  enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateIngredientCategoryDTO}) => {
    return ingredientsService.updateIngredientCategory(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories-details', categoryId] })
		router.push({ name: getRouteName("ADMIN_INGREDIENT_CATEGORIES") })
	},
})

function handleUpdate(updatedData: UpdateIngredientCategoryDTO) {
  if (isNaN(Number(categoryId))) return router.back()

	updateMutation.mutate({id: Number(categoryId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
