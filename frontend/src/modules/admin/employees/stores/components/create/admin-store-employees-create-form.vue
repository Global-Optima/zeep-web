<template>
	<AdminEmployeesCreateForm
		title="Создать сотрудника кафе"
		formTitle="Создайте сотрудника кафе"
		formDescription="Заполните форму ниже, чтобы создать данные сотрудника."
		:roles="storeRoles"
		@onSubmit="handleCreateStoreEmployee"
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

// Define store-specific roles
const storeRoles = ref([
  { value: EmployeeRole.STORE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.BARISTA, label: 'Бариста' }
])

const handleCreateStoreEmployee = (formValues: CreateEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
