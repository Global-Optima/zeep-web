<template>
	<AdminAdditiveCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminAdditiveCreateForm from '@/modules/admin/additives/components/create/admin-additive-create-form.vue'
import type { CreateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateAdditiveDTO) => additivesService.createAdditive(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
		router.back()
	},
})

function handleCreate(dto: CreateAdditiveDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
