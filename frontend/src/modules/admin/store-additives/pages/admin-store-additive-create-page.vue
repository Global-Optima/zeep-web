<template>
	<AdminStoreAdditiveCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreAdditiveCreateForm from '@/modules/admin/store-additives/components/create/admin-store-additive-create-form.vue'
import type { CreateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateStoreAdditiveDTO[]) => storeAdditivesService.createStoreAdditive(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-additives'] })
		router.back()
	},
})

function handleCreate(dto: CreateStoreAdditiveDTO[]) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
