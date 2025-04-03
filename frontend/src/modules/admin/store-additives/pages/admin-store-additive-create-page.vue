<template>
	<AdminStoreAdditiveCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreAdditiveCreateForm from '@/modules/admin/store-additives/components/create/admin-store-additive-create-form.vue'
import type { CreateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStoreAdditiveDTO[]) => storeAdditivesService.createStoreAdditive(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Добавление новых модификаторов кафе. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-additives'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Модификаторы кафе успешно добавлены.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении модификаторов кафе.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateStoreAdditiveDTO[]) {
	if (dto.length === 0) {
		toast({
			title: 'Ошибка',
			description: 'Список модификаторов пуст. Пожалуйста, добавьте модификатора перед сохранением.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
