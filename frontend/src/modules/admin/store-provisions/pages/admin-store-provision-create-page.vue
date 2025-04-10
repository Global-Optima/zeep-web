<template>
	<AdminStoreProvisionCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import AdminStoreProvisionCreateForm from '@/modules/admin/store-provisions/components/create/admin-store-provision-create-form..vue'
import type { CreateStoreProvisionDTO } from '@/modules/admin/store-provisions/models/store-provision.models'
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const {toastLocalizedError} = useAxiosLocaleToast()

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
	onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, 'Произошла ошибка при создании заготовки.')
	},
})

function handleCreate(dto: CreateStoreProvisionDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
