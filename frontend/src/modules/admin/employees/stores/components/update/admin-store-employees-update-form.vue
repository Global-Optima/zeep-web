<template>
	<AdminEmployeesUpdateForm
		:employee="employee"
		title="Обновить сотрудника кафе"
		:roles="storeRoles"
		@onSubmit="handleUpdateStoreEmployee"
		@onCancel="handleCancel"
	/>
</template>

<script setup lang="ts">
import AdminEmployeesUpdateForm from '@/modules/admin/employees/components/update/admin-employees-update-form.vue'
import { EmployeeRole, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type { UpdateStoreEmployeeDTO } from '@/modules/admin/employees/stores/models/store-employees.model'
import { ref } from 'vue'

defineProps<{
  employee: EmployeeDTO
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateStoreEmployeeDTO): void
  (e: 'onCancel'): void
}>()

const storeRoles = ref([
  { value: EmployeeRole.STORE_MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.BARISTA, label: 'Бариста' }
])

const handleUpdateStoreEmployee = (formValues: UpdateStoreEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
