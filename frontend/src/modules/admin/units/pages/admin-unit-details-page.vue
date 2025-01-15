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
import { getRouteName } from '@/core/config/routes.config'
import AdminUnitsDetailsForm from '@/modules/admin/units/components/details/admin-units-details-form.vue'
import type { UpdateUnitDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const unitId = route.params.id as string

const queryClient = useQueryClient()

const { data: unitDetails } = useQuery({
  queryKey: computed(() => ['admin-unit-details', unitId]),
	queryFn: () => unitsService.getUnitByID(Number(unitId)),
  enabled: !isNaN(Number(unitId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateUnitDTO}) => {
    return unitsService.updateUnit(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-units'] })
		queryClient.invalidateQueries({ queryKey: ['admin-unit-details', unitId] })
		router.push({ name: getRouteName("ADMIN_UNITS") })
	},
})

function handleUpdate(updatedData: UpdateUnitDTO) {
  if (isNaN(Number(unitId))) return router.back()

	updateMutation.mutate({id: Number(unitId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
