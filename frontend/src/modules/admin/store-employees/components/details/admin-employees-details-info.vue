<template>
	<div class="relative bg-white p-6 border rounded-xl">
		<!-- Employee Image -->
		<!-- TODO: Add image url here -->
		<img
			class="mb-6 rounded-full w-48 h-48 object-cover"
			src=""
			alt="Employee Image"
		/>
		<!-- Employee Attributes -->
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
		<!-- Action Buttons -->
		<div class="top-6 right-6 absolute">
			<Button
				size="icon"
				variant="ghost"
			>
				<Pencil class="w-6 h-6 text-gray-500" />
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import Button from '@/core/components/ui/button/Button.vue'
import { EMPLOYEE_ROLES_FORMATTED, type Employee } from '@/modules/admin/employees/models/employees.models'
import { Pencil } from 'lucide-vue-next'
import { computed } from 'vue'

const {employee} = defineProps<{employee: Employee}>()

const attributes = computed(() => [
  { label: 'Имя', value: `${employee.firstName} ${employee.lastName}` },
  { label: 'Должность', value: EMPLOYEE_ROLES_FORMATTED[employee.role] },
  { label: 'Email', value: employee.email },
  { label: 'Телефон', value: employee.phone },
  { label: 'Статус', value: employee.isActive ? "Активен" : "Деактивирован" },
]);
</script>

<style scoped></style>
