<template>
	<div class="relative bg-white p-6 border rounded-xl">
		<div class="flex justify-between items-start gap-4">
			<Avatar class="bg-gray-200 rounded-lg size-32">
				<AvatarFallback
					>{{ employee.firstName.charAt(0) }}{{ employee.lastName.charAt(0) }}</AvatarFallback
				>
			</Avatar>

			<div class="flex items-center gap-1">
				<Button
					v-if="showReassignButton"
					size="icon"
					variant="ghost"
					@click="onReassignEmployeeClick"
				>
					<ArrowRightLeft class="size-5 text-gray-500" />
				</Button>
				<Button
					size="icon"
					variant="ghost"
					@click="onUpdateEmployeeClick"
				>
					<Pencil class="size-5 text-gray-500" />
				</Button>
			</div>
		</div>

		<div
			v-for="attribute in attributes"
			:key="attribute.label"
			class="mt-4"
		>
			<label class="block text-gray-400 text-sm">
				{{ attribute.label }}
			</label>
			<p class="mt-1 text-gray-900 text-base">{{ attribute.value }}</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/core/components/ui/avatar'
import Button from '@/core/components/ui/button/Button.vue'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EMPLOYEE_ROLES_FORMATTED, EmployeeRole, type BaseEmployeeDetailsDTO } from '@/modules/admin/employees/models/employees.models'
import { ArrowRightLeft, Pencil } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const {employee} = defineProps<{employee: BaseEmployeeDetailsDTO & {id: number}}>()

const router = useRouter()

const showReassignButton = useHasRole([EmployeeRole.ADMIN])

const onUpdateEmployeeClick = () => {
  router.push(`/admin/employees/${employee.type.toLowerCase()}/${employee.id}/update`)
}

const onReassignEmployeeClick = () => {
  router.push(`/admin/employees/${employee.employeeId}/reassign`)
}

const attributes = computed(() => [
  { label: 'Имя', value: `${employee.firstName} ${employee.lastName}` },
  { label: 'Должность', value: EMPLOYEE_ROLES_FORMATTED[employee.role] },
  { label: 'Email', value: employee.email },
  { label: 'Телефон', value: employee.phone },
  { label: 'Статус', value: employee.isActive ? "Активен" : "Деактивирован" },
]);
</script>

<style scoped></style>
