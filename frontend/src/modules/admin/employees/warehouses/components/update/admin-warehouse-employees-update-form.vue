<template>
	<AdminEmployeesUpdateForm
		:employee="employee"
		title="Обновить сотрудника склада"
		:roles="warehouseRoles"
		@onSubmit="handleUpdateWarehouseEmployee"
		@onCancel="handleCancel"
	/>
</template>

<script setup lang="ts">
import AdminEmployeesUpdateForm from '@/modules/admin/employees/components/update/admin-employees-update-form.vue'
import { EmployeeRole, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type { UpdateWarehouseEmployeeDTO } from '@/modules/admin/employees/warehouses/models/warehouse-employees.model'
import { ref } from 'vue'

defineProps<{
  employee: EmployeeDTO
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateWarehouseEmployeeDTO): void
  (e: 'onCancel'): void
}>()

// Define warehouse-specific roles
const warehouseRoles = ref([
  { value: EmployeeRole.WAREHOUSE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.WAREHOUSE_EMPLOYEE, label: 'Сотрудник склада' }
])

const handleUpdateWarehouseEmployee = (formValues: UpdateWarehouseEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
