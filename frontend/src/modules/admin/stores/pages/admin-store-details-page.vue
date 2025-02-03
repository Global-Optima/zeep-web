<template>
	<AdminStoreDetailsForm
		v-if="storeData"
		:initialData="storeData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
	<div v-else>Загрузка...</div>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreDetailsForm from '@/modules/admin/stores/components/details/admin-store-details-form.vue'
import type { UpdateStoreDTO } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()

const storeId = route.params.id as string

const queryClient = useQueryClient()

const { data: storeData } = useQuery({
	queryKey: ['store', storeId],
	queryFn: () => storesService.getStore(Number(storeId)),
	enabled: computed(() => !!storeId),
})

const updateMutation = useMutation({
	mutationFn: (updatedData: UpdateStoreDTO) =>
		storesService.updateStore(Number(storeId), updatedData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Пожалуйста, подождите, данные обновляются.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		queryClient.invalidateQueries({ queryKey: ['store', storeId] })
		toast({
			title: 'Успех!',
			description: 'Данные магазина успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных магазина.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStoreDTO) {
	updateMutation.mutate(updatedData)
}

function handleCancel() {
	router.back()
}
</script>
