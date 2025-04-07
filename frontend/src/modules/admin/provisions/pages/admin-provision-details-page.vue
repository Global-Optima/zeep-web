<template>
	<p v-if="!ingredientDetails">Заготовка не найдена</p>

	<AdminProvisionUpdateForm
		v-else
		:provision="ingredientDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProvisionUpdateForm from "@/modules/admin/provisions/components/details/admin-provision-update-form.vue"
import type { UpdateProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const provisionId = route.params.id as string

const { data: ingredientDetails } = useQuery({
	queryKey: computed(() => ['admin-provision-details', provisionId]),
	queryFn: () => provisionsService.getProvisionById(Number(provisionId)),
	enabled: !isNaN(Number(provisionId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateProvisionDTO }) =>
    provisionsService.updateProvision(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-provisions'] })
		queryClient.invalidateQueries({ queryKey: ['admin-provision-details', provisionId] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Данные сырья успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateProvisionDTO) {
	if (isNaN(Number(provisionId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(provisionId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
