<template>
	<AdminEmployeesUpdateForm
		:employee="employee"
		title="Обновить сотрудника франчайзи"
		:roles="franchiseeRoles"
		@onSubmit="handleUpdateFranchiseeEmployee"
		@onCancel="handleCancel"
	/>
</template>

<script setup lang="ts">
import AdminEmployeesUpdateForm from '@/modules/admin/employees/components/update/admin-employees-update-form.vue'
import type { UpdateFranchiseeEmployeeDTO } from '@/modules/admin/employees/franchisees/models/franchisees-employees.model'
import { EmployeeRole, type BaseEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { ref } from 'vue'

defineProps<{
  employee: BaseEmployeeDTO
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateFranchiseeEmployeeDTO): void
  (e: 'onCancel'): void
}>()

// Define franchisee-specific roles
const franchiseeRoles = ref([
  { value: EmployeeRole.FRANCHISEE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.FRANCHISEE_OWNER, label: 'Владелец' }
])

const handleUpdateFranchiseeEmployee = (formValues: UpdateFranchiseeEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
