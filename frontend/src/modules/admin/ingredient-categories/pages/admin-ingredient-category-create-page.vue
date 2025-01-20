<template>
	<AdminIngredientCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminIngredientCategoryCreateForm from '@/modules/admin/ingredient-categories/components/create/admin-ingredient-category-create-form.vue'
import type { CreateIngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateIngredientCategoryDTO) => ingredientsService.createIngredientCategory(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories'] })
		router.back()
	},
})

function handleCreate(dto: CreateIngredientCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
