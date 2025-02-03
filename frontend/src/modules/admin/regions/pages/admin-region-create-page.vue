<template>
	<AdminRegionCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminRegionCreateForm from '@/modules/admin/regions/components/create/admin-region-create-form.vue'
import type { CreateRegionDTO } from '@/modules/admin/regions/models/regions.model'
import { regionsService } from '@/modules/admin/regions/services/regions.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateRegionDTO) => regionsService.create(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового региона. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-regions'] })
		toast({
			title: 'Успех!',
			description: 'Регион успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании региона.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateRegionDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
