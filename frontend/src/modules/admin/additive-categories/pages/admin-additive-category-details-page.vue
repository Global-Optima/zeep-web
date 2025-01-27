<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminAdditiveCategoryDetailsForm
		v-else
		:category="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminAdditiveCategoryDetailsForm from '@/modules/admin/additive-categories/components/details/admin-additive-category-details-form.vue'
import type { UpdateAdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const categoryId = route.params.id as string

const { data: categoryDetails } = useQuery({
	queryKey: computed(() => ['admin-additive-categories-details', categoryId]),
	queryFn: () => additivesService.getAdditiveCategoryById(Number(categoryId)),
	enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateAdditiveCategoryDTO }) =>
		additivesService.updateAdditiveCategory(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных категории добавок. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories-details', categoryId] })
		toast({
			title: 'Успех!',
			description: 'Данные категории успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении категории добавок.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateAdditiveCategoryDTO) {
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
