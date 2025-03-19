<template>
	<p v-if="!unitDetails">Размер не найден</p>

	<AdminUnitsDetailsForm
		v-else
		:unit="unitDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminUnitsDetailsForm from '@/modules/admin/units/components/details/admin-units-details-form.vue'
import type { UpdateUnitDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()

const unitId = route.params.id as string

const queryClient = useQueryClient()

const { data: unitDetails } = useQuery({
	queryKey: computed(() => ['admin-unit-details', unitId]),
	queryFn: () => unitsService.getUnitByID(Number(unitId)),
	enabled: !isNaN(Number(unitId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateUnitDTO }) => unitsService.updateUnit(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Пожалуйста, подождите, данные обновляются.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-units'] })
		queryClient.invalidateQueries({ queryKey: ['admin-unit-details', unitId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Размер успешно обновлён.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении размера.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateUnitDTO) {
	if (isNaN(Number(unitId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор размера.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(unitId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
