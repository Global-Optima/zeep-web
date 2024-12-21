<template>
	<AdminSuppliersCreateForm
		:initialData="defaultStoreData"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminSuppliersCreateForm from '@/modules/admin/suppliers/components/create/admin-suppliers-create-form.vue'
import type { CreateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const defaultStoreData: Partial<CreateSupplierDTO> = {
	name: '',
  address: '',
	contactPhone: '',
	contactEmail: '',
}

const createMutation = useMutation({
	mutationFn: (newStoreData: CreateSupplierDTO) => suppliersService.createSupplier(newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['suppliers'] })
		router.push({ name: getRouteName("ADMIN_SUPPLIERS") })
	},
})

function handleCreate(dto: CreateSupplierDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
