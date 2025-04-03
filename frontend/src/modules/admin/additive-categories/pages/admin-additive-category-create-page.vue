<template>
	<AdminAdditiveCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminAdditiveCategoryCreateForm from '@/modules/admin/additive-categories/components/create/admin-additive-category-create-form.vue'
import type { CreateAdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateAdditiveCategoryDTO) => additivesService.createAdditiveCategory(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой категории модификаторов. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Категория модификаторов успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании категории модификаторов.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateAdditiveCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
