<template>
	<AdminStoreCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreCreateForm from '@/modules/admin/stores/components/create/admin-store-create-form.vue'
import type { CreateStoreDTO } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (newStoreData: CreateStoreDTO) => storesService.createStore(newStoreData),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Пожалуйста, подождите, создается новый кафе.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Кафе успешно создан.',
		})
		router.push({ name: 'ADMIN_STORES' })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании кафе.',
			variant: 'destructive',
		})
	},
})

function handleCreate(newStoreData: CreateStoreDTO) {
	createMutation.mutate(newStoreData)
}

function handleCancel() {
	router.back()
}
</script>
