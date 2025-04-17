<template>
	<p v-if="!ingredientDetails">Сырье не найдено</p>

	<AdminIngredientsDetailsForm
		v-else
		:ingredient="ingredientDetails"
		:readonly="!canUpdate"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminIngredientsDetailsForm from '@/modules/admin/ingredients/components/details/admin-ingredients-details-form.vue'
import type { UpdateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const canUpdate = useHasRole([EmployeeRole.ADMIN])

const ingredientId = route.params.id as string

const { data: ingredientDetails } = useQuery({
	queryKey: computed(() => ['admin-ingredient-details', ingredientId]),
	queryFn: () => ingredientsService.getIngredientById(Number(ingredientId)),
	enabled: !isNaN(Number(ingredientId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateIngredientDTO }) =>
		ingredientsService.updateIngredient(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных сырья. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-ingredients'] })
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-details', ingredientId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные сырья успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении сырья.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateIngredientDTO) {
	if (isNaN(Number(ingredientId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор сырья.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(ingredientId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
