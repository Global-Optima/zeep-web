<template>
	<p v-if="!storeProvisionDTO || !provisionDTO">Заготовка не найдена</p>

	<AdminStoreProvisionUpdateForm
		v-else
		:store-provision="storeProvisionDTO"
		:provision="provisionDTO"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import { provisionsService } from '@/modules/admin/provisions/services/provisions.service'
import AdminStoreProvisionUpdateForm from '@/modules/admin/store-provisions/components/update/admin-store-provision-update-form.vue'
import type { UpdateStoreProvisionDTO } from '@/modules/admin/store-provisions/models/store-provision.models'
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const {toastLocalizedError} = useAxiosLocaleToast()

const storeProvisionId = route.params.id as string

const { data: storeProvisionDTO } = useQuery({
	queryKey: computed(() => ['admin-store-provision-details', storeProvisionId]),
	queryFn: () => storeProvisionsService.getStoreProvisionById(Number(storeProvisionId)),
	enabled: !isNaN(Number(storeProvisionId)),
})

const { data: provisionDTO } = useQuery({
	queryKey: computed(() => ['admin-provision-details', storeProvisionDTO.value?.provision.id.toString()]),
	queryFn: () => provisionsService.getProvisionById(Number(storeProvisionDTO.value?.provision.id)),
	enabled: computed(() => Boolean(storeProvisionDTO.value?.provision)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStoreProvisionDTO }) =>
    storeProvisionsService.updateStoreProvision(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-provisions'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-provision-details', storeProvisionId] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Данные успешно обновлены.',
		})
	},
	onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, 'Произошла ошибка при обновлении.')
	},
})

function handleUpdate(updatedData: UpdateStoreProvisionDTO) {
	if (isNaN(Number(storeProvisionId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(storeProvisionId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
