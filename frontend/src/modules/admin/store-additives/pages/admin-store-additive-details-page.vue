<template>
	<p v-if="!storeAdditiveDetails">Модификатор не найден</p>

	<AdminStoreAdditiveDetailsForm
		v-else
		:initialAdditive="storeAdditiveDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
		:readonly="!canUpdate"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminStoreAdditiveDetailsForm from '@/modules/admin/store-additives/components/details/admin-store-additive-details-form.vue'
import type { UpdateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const storeAdditiveId = route.params.id as string

const { data: storeAdditiveDetails } = useQuery({
	queryKey: computed(() => ['admin-store-additive-details', storeAdditiveId]),
	queryFn: () => storeAdditivesService.getStoreAdditiveById(Number(storeAdditiveId)),
	enabled: !isNaN(Number(storeAdditiveId)),
})

const canUpdate = useHasRole([EmployeeRole.STORE_MANAGER])

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStoreAdditiveDTO }) =>
		storeAdditivesService.updateStoreAdditive(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных модификатора кафе. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-additives'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-additive-details', storeAdditiveId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные модификатора кафе успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных модификатора кафе.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStoreAdditiveDTO) {
	if (isNaN(Number(storeAdditiveId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор модификатора.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(storeAdditiveId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
