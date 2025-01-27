<template>
	<AdminEmployeesUpdateForm
		v-if="employee"
		:employee="employee"
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>

	<p v-else>Сотрудник не найден</p>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminEmployeesUpdateForm from '@/modules/admin/store-employees/components/update/admin-employees-update-form.vue'
import type { UpdateEmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
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
	mutationFn: (newStoreData: UpdateEmployeeDTO) => employeesService.updateStoreEmployee(Number(employeeId), newStoreData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-employees'] })
    queryClient.invalidateQueries({ queryKey: ['store-employee', employeeId] })
		router.push({ name: getRouteName("ADMIN_STORE_EMPLOYEES") })
	},
})

function handleCreate(dto: UpdateEmployeeDTO) {
	updateMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
