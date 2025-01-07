<template>
	<AdminIngredientCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminIngredientCreateForm from '@/modules/admin/ingredients/components/create/admin-ingredient-create-form.vue'
import type { CreateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateIngredientDTO) => ingredientsService.createIngredient(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredients'] })
		router.back()
	},
})

function handleCreate(dto: CreateIngredientDTO) {
  const formattedDateDTO = {...dto, expiresAt: dto.expiresAt ? new Date(dto.expiresAt).toISOString() : new Date().toISOString()}

	createMutation.mutate(formattedDateDTO)
}

function handleCancel() {
	router.back()
}
</script>
