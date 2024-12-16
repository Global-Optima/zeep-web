<template>
	<AdminStoreCreateForm
		:initialData="defaultStoreData"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreCreateForm from '@/modules/admin/stores/components/create/admin-store-create-form.vue'
import type { CreateStoreDTO } from '@/modules/stores/models/stores-dto.model'
import { storesService } from '@/modules/stores/services/stores.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

// Default values for a new store
const defaultStoreData: Partial<CreateStoreDTO> = {
	name: '',
	isFranchise: false,
	facilityAddress: {
		address: '',
		longitude: 0,
		latitude: 0,
	},
	contactPhone: '',
	contactEmail: '',
	storeHours: '',
}

// Mutation for creating a store
const createMutation = useMutation({
	mutationFn: (newStoreData: CreateStoreDTO) => storesService.createStore(newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		router.push({ name: 'ADMIN_STORES' })
	},
})

function handleCreate(newStoreData: CreateStoreDTO) {
  console.log("HEREEEE")
	createMutation.mutate(newStoreData)
}

function handleCancel() {
	router.back()
}
</script>
