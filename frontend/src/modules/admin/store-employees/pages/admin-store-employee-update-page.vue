<template>
	<AdminEmployeesUpdateForm
		v-if="employee"
		:initialData="employee"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminEmployeesUpdateForm from '@/modules/admin/employees/components/update/admin-employees-update-form.vue'
import type { UpdateEmployeeDto } from '@/modules/admin/employees/models/employees.models'
import { employeesService } from '@/modules/admin/employees/services/employees.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
  queryKey: ['employee', employeeId],
	queryFn: () => employeesService.getStoreEmployeeById(Number(employeeId)),
  enabled: !!employeeId,
})

const updateMutation = useMutation({
	mutationFn: (newStoreData: UpdateEmployeeDto) => employeesService.updateStoreEmployee(Number(employeeId), newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['employees'] })
    queryClient.invalidateQueries({ queryKey: ['employee', employeeId] })
		router.push({ name: getRouteName("ADMIN_EMPLOYEES") })
	},
})

function handleCreate(dto: UpdateEmployeeDto) {
	updateMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
