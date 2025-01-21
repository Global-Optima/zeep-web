<template>
	<AdminSupplierCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminSupplierCreateForm from '@/modules/admin/suppliers/components/create/admin-supplier-create-form.vue'
import type { CreateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateSupplierDTO) => suppliersService.createSupplier(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
		router.back()
	},
})

function handleCreate(dto: CreateSupplierDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
