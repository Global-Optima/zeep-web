<template>
	<p v-if="!additiveDetails">Топпинг не найден</p>

	<AdminAdditiveDetailsForm
		v-else
    ref="formRef"
		:additive="additiveDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
		:readonly="!canUpdateAdditives"
		:isSubmitting="isPending"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import AdminAdditiveDetailsForm from '@/modules/admin/additives/components/details/admin-additive-details-form.vue'
import type { UpdateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'
import {useTemplateRef} from "vue";

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const canUpdateAdditives = useHasRole([EmployeeRole.ADMIN])

const additiveId = route.params.id as string

const { data: additiveDetails } = useQuery({
	queryKey: ['admin-additive-details', additiveId],
	queryFn: () => additivesService.getAdditiveById(Number(additiveId)),
	enabled: !isNaN(Number(additiveId)),
})

const formRef = useTemplateRef<InstanceType<typeof AdminAdditiveDetailsForm>>('formRef')

const {mutate, isPending} = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateAdditiveDTO }) =>
		additivesService.updateAdditive(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных топпинга. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
		queryClient.invalidateQueries({ queryKey: ['admin-additive-details', additiveId] })

    formRef.value?.resetFormValues();

		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Данные топпинга успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных топпинга.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateAdditiveDTO) {
	if (isNaN(Number(additiveId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор топпинга.',
			variant: 'destructive',
		})
		return router.back()
	}

	mutate({ id: Number(additiveId), dto: data })
}

function handleCancel() {
	router.back()
}
</script>
