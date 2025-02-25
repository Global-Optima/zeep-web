<template>
	<AdminEmployeesCreateForm
		title="Создать администратора"
		formTitle="Создайте администратора"
		formDescription="Заполните форму ниже, чтобы создать учетную запись администратора."
		:roles="adminRoles"
		@onSubmit="handleCreateAdmin"
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

const adminRoles = ref([
  { value: EmployeeRole.ADMIN, label: 'Администратор' }
])

const handleCreateAdmin = (formValues: CreateEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
