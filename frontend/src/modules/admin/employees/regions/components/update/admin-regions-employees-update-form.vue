<template>
	<AdminEmployeesUpdateForm
		:employee="employee"
		title="Обновить сотрудника региона"
		:roles="regionRoles"
		@onSubmit="handleUpdateRegionEmployee"
		@onCancel="handleCancel"
	/>
</template>

<script setup lang="ts">
import AdminEmployeesUpdateForm from '@/modules/admin/employees/components/update/admin-employees-update-form.vue'
import { EmployeeRole, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type { UpdateRegionEmployeeDTO } from '@/modules/admin/employees/regions/models/region-employees.model'
import { ref } from 'vue'

defineProps<{
  employee: EmployeeDTO
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateRegionEmployeeDTO): void
  (e: 'onCancel'): void
}>()

const regionRoles = ref([
  { value: EmployeeRole.REGION_WAREHOUSE_MANAGER, label: 'Менеджер' }
])

const handleUpdateRegionEmployee = (formValues: UpdateRegionEmployeeDTO) => {
  emit('onSubmit', formValues)
}

const handleCancel = () => {
  emit('onCancel')
}
</script>
