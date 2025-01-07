<template>
	<p v-if="!ingredientDetails">Ингредиент не найден</p>

	<AdminIngredientsDetailsForm
		v-else
		:ingredient="ingredientDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminIngredientsDetailsForm from '@/modules/admin/ingredients/components/details/admin-ingredients-details-form.vue'
import type { UpdateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const ingredientId = route.params.id as string

const queryClient = useQueryClient()


const { data: ingredientDetails } = useQuery({
  queryKey: computed(() => ['admin-ingredient-details', ingredientId]),
	queryFn: () => ingredientsService.getIngredientById(Number(ingredientId)),
  enabled: !isNaN(Number(ingredientId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateIngredientDTO}) => {
    return ingredientsService.updateIngredient(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredients'] })
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-details', ingredientId] })
		router.push({ name: getRouteName("ADMIN_INGREDIENTS") })
	},
})

function handleUpdate(updatedData: UpdateIngredientDTO) {
  if (isNaN(Number(ingredientId))) return router.back()

  const dto: UpdateIngredientDTO = {...updatedData, expiresAt: updatedData.expiresAt ? new Date(updatedData.expiresAt).toISOString() : new Date().toISOString()}

	updateMutation.mutate({id: Number(ingredientId), dto})
}

function handleCancel() {
	router.back()
}
</script>
