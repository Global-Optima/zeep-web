<template>
	<AdminStoreProvisionCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreProvisionCreateForm from '@/modules/admin/store-provisions/components/create/admin-store-provision-create-form..vue'
import type { CreateStoreProvisionDTO } from '@/modules/admin/store-provisions/models/store-provision.models'
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStoreProvisionDTO) => storeProvisionsService.createStoreProvision(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-provisions'] })
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

function handleCreate(dto: CreateStoreProvisionDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
