<template>
	<AdminProvisionCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProvisionCreateForm from "@/modules/admin/provisions/components/create/admin-provision-create-form.vue"
import type { CreateProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateProvisionDTO) => provisionsService.createProvision(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-provisions'] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Заготовка успешна создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании заготовки.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateProvisionDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
