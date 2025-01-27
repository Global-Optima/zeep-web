<template>
	<AdminAdditiveCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminAdditiveCreateForm from '@/modules/admin/additives/components/create/admin-additive-create-form.vue'
import type { CreateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateAdditiveDTO) => additivesService.createAdditive(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой добавки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
		toast({
			title: 'Успех!',
			description: 'Добавка успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании добавки.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateAdditiveDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
