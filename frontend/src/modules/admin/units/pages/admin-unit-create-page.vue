<template>
	<AdminUnitsCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminUnitsCreateForm from '@/modules/admin/units/components/create/admin-units-create-form.vue'
import type { CreateUnitDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateUnitDTO) => unitsService.createUnit(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-units'] })
		router.back()
	},
})

function handleCreate(dto: CreateUnitDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
