<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Имя</TableHead>
				<TableHead class="p-4">Роль</TableHead>
				<TableHead class="p-4">Телефон</TableHead>
				<TableHead class="p-4">Статус</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="employee in employees"
				:key="employee.id"
				class="h-12 cursor-pointer"
				@click="goToEmployee(employee.id)"
			>
				<!-- Employee Name and Image -->
				<TableCell class="p-4">
					<span class="font-medium"> {{ employee.firstName }} {{ employee.lastName }}</span>
				</TableCell>
				<!-- Role -->
				<TableCell class="p-4">
					{{ EMPLOYEE_ROLES_FORMATTED[employee.role] }}
				</TableCell>

				<TableCell class="p-4">
					{{ formatPhoneNumber(employee.phone) }}
				</TableCell>

				<TableCell class="p-4">
					<span
						:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs ',
                STATUS_COLOR[employee.isActive ? 'active' : 'disabled'],
              ]"
					>
						{{ STATUS_FORMATTED[employee.isActive ? "active" : "disabled"] }}
					</span>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { formatPhoneNumber } from '@/core/utils/fomat-phone-number.utils'
import type { FranchiseeEmployeeDTO } from '@/modules/admin/employees/franchisees/models/franchisees-employees.model'
import { EMPLOYEE_ROLES_FORMATTED } from '@/modules/admin/employees/models/employees.models'

const {employees} = defineProps<{employees: FranchiseeEmployeeDTO[]}>()

const router = useRouter()

// Navigate to employee details
const goToEmployee = (employeeId: number) => {
  router.push(`/admin/employees/franchisee/${employeeId}`)
}

// Status colors and formatted text
const STATUS_COLOR: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  disabled: 'bg-red-100 text-red-800',
}

const STATUS_FORMATTED: Record<string, string> = {
  active: 'Активен',
  disabled: 'Неактивный',
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
