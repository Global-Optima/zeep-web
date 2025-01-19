<template>
	<AdminStoreEmployeesCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreEmployeesCreateForm from '@/modules/admin/store-employees/components/create/admin-store-employees-create-form.vue'
import { type CreateStoreEmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (newStoreData: CreateStoreEmployeeDTO) => employeesService.createStoreEmployee(newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-employees'] })
		router.push({ name: getRouteName("ADMIN_STORE_EMPLOYEES") })
	},
})

function handleCreate(dto: CreateStoreEmployeeDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
