<template>
	<div class="relative bg-white p-6 border rounded-xl">
		<Avatar class="bg-gray-200 rounded-lg size-32">
			<AvatarFallback>{{ getEmployeeInitials(employee) }}</AvatarFallback>
		</Avatar>
		<div
			v-for="attribute in attributes"
			:key="attribute.label"
			class="mt-4"
		>
			<label class="block text-gray-400 text-sm">
				{{ attribute.label }}
			</label>
			<p class="mt-1 text-base text-gray-900">{{ attribute.value }}</p>
		</div>
		<div class="top-6 right-6 absolute">
			<Button
				size="icon"
				variant="ghost"
				@click="$router.push(`/admin/employees/store/${employee.id}/update`)"
			>
				<Pencil class="w-6 h-6 text-gray-500" />
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/core/components/ui/avatar'
import Button from '@/core/components/ui/button/Button.vue'
import { getEmployeeInitials } from '@/core/utils/user-formatting.utils'
import { EMPLOYEE_ROLES_FORMATTED, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { Pencil } from 'lucide-vue-next'
import { computed } from 'vue'

const {employee} = defineProps<{employee: EmployeeDTO}>()

const attributes = computed(() => [
  { label: 'Имя', value: `${employee.firstName} ${employee.lastName}` },
  { label: 'Должность', value: EMPLOYEE_ROLES_FORMATTED[employee.role] },
  { label: 'Email', value: employee.email },
  { label: 'Телефон', value: employee.phone },
  { label: 'Статус', value: employee.isActive ? "Активен" : "Деактивирован" },
]);
</script>

<style scoped></style>
