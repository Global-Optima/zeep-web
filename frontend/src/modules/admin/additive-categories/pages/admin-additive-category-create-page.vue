<template>
	<AdminAdditiveCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminAdditiveCategoryCreateForm from '@/modules/admin/additive-categories/components/create/admin-additive-category-create-form.vue'
import type { CreateAdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateAdditiveCategoryDTO) => additivesService.createAdditiveCategory(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories'] })
		router.back()
	},
})

function handleCreate(dto: CreateAdditiveCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
