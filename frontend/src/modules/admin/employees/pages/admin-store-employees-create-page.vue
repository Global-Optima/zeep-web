<template>
	<AdminEmployeesCreateFormStore
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminEmployeesCreateFormStore from '@/modules/admin/employees/components/create/admin-employees-create-form-store.vue'
import { EmployeeType, type CreateEmployeeDto } from '@/modules/employees/models/employees.models'
import { employeesService } from '@/modules/employees/services/employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (newStoreData: CreateEmployeeDto) => employeesService.createEmployee(newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['employees'] })
		router.push({ name: getRouteName("ADMIN_EMPLOYEES") })
	},
})

function handleCreate(dto: CreateEmployeeDto) {
  dto.type = EmployeeType.STORE
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
