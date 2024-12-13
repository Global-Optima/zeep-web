<template>
	<AdminStoreManage
		:isEditing="false"
		:initialData="defaultStoreData"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreManage from '@/modules/admin/stores/components/details/admin-store-manage.vue'
import type { Store } from '@/modules/stores/models/stores.models'
import { storesService } from '@/modules/stores/services/stores.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

// Default values for a new store
const defaultStoreData: Partial<Store> = {
	name: '',
	isFranchise: false,
	facilityAddress: {
		id: 0,
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
	mutationFn: (newStoreData: Store) => storesService.createStore(newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		router.push({ name: 'ADMIN_STORES' })
	},
})

function handleCreate(newStoreData: Store) {
	createMutation.mutate(newStoreData)
}

function handleCancel() {
	router.back()
}
</script>
