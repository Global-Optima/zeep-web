<template>
	<AdminEmployeesCreateForm
		title="Создать сотрудника склада"
		formTitle="Создайте сотрудника склада"
		formDescription="Заполните форму ниже, чтобы создать данные сотрудника."
		:roles="warehouseRoles"
		@onSubmit="handleCreateWarehouseEmployee"
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

// Define warehouse-specific roles
const warehouseRoles = ref([
  { value: EmployeeRole.WAREHOUSE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.WAREHOUSE_EMPLOYEE, label: 'Сотрудник склада' }
])

const handleCreateWarehouseEmployee = (formValues: CreateEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
