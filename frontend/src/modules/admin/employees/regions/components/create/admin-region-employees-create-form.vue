<template>
	<AdminEmployeesCreateForm
		title="Создать сотрудника региона"
		formTitle="Создайте сотрудника региона"
		formDescription="Заполните форму ниже, чтобы создать данные сотрудника."
		:roles="regionRoles"
		@onSubmit="handleCreateRegionEmployee"
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

// Define region-specific roles
const regionRoles = ref([
  { value: EmployeeRole.REGION_WAREHOUSE_MANAGER, label: 'Менеджер' }
])

const handleCreateRegionEmployee = (formValues: CreateEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
