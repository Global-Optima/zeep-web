<template>
	<AdminEmployeesCreateForm
		title="Создать сотрудника франчайзи"
		formTitle="Создайте сотрудника франчайзи"
		formDescription="Заполните форму ниже, чтобы создать данные сотрудника."
		:roles="franchiseeRoles"
		@onSubmit="handleCreateFranchiseeEmployee"
		@onCancel="handleCancel"
	/>
</template>

<script setup lang="ts">
import AdminEmployeesCreateForm from '@/modules/admin/employees/components/create/admin-employees-create-form.vue'
import { EmployeeRole, type CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'onSubmit', formValues: CreateEmployeeDTO): void
  (e: 'onCancel'): void
}>()

// Define Franchisee-specific roles
const franchiseeRoles = ref([
  { value: EmployeeRole.FRANCHISEE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.FRANCHISEE_OWNER, label: 'Владелец' }
])

const handleCreateFranchiseeEmployee = (formValues: CreateEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
