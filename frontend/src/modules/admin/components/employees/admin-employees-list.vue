<template>
	<div>
		<!-- If no employees, display message -->
		<p
			v-if="employees.length === 0"
			class="text-muted-foreground"
		>
			Сотрудники не найдены
		</p>
		<!-- If there are employees, display the table -->
		<Table v-else>
			<TableHeader>
				<TableRow>
					<TableHead>Сотрудник</TableHead>
					<TableHead>Роль</TableHead>
					<TableHead>Телефон</TableHead>
					<TableHead>Магазин</TableHead>
					<TableHead>Статус</TableHead>
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
					<TableCell class="flex items-center space-x-3">
						<img
							:src="employee.avatar"
							alt="Avatar"
							class="rounded-full w-8 h-8 object-cover"
						/>
						<span class="font-medium"> {{ employee.firstName }} {{ employee.lastName }} </span>
					</TableCell>
					<!-- Role -->
					<TableCell class="font-medium">
						<span
							:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
              ]"
						>
							{{ ROLE_FORMATTED[employee.role] }}
						</span>
					</TableCell>
					<!-- Phone Number -->
					<TableCell class="font-medium">
						{{ formatPhoneNumber(employee.phone) }}
					</TableCell>
					<!-- Working Store Name -->
					<TableCell class="font-medium">
						{{ employee.storeName }}
					</TableCell>
					<!-- Status -->
					<TableCell class="font-medium">
						<span
							:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
                STATUS_COLOR[employee.status],
              ]"
						>
							{{ STATUS_FORMATTED[employee.status] }}
						</span>
					</TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'

// Mock data for employees (in Russian)
interface Employee {
  id: number
  firstName: string
  lastName: string
  avatar: string
  role: string
  phone: string
  storeName: string
  status: string
}

const employees = ref<Employee[]>([
  {
    id: 1,
    firstName: 'Алексей',
    lastName: 'Смирнов',
    avatar: 'https://via.placeholder.com/32x32', // Placeholder image
    role: 'manager',
    phone: '+79161234567',
    storeName: 'Магазин №1',
    status: 'active',
  },
  {
    id: 2,
    firstName: 'Екатерина',
    lastName: 'Иванова',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'cashier',
    phone: '+79169876543',
    storeName: 'Магазин №2',
    status: 'disabled',
  },
  // Add more mock employees as needed
])

const router = useRouter()

// Navigate to employee details
const goToEmployee = (employeeId: number) => {
  router.push(`/admin/employees/${employeeId}`)
}

// Format phone number for display
const formatPhoneNumber = (phone: string) => {
  // Simple formatting, adjust as needed
  return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5')
}

const ROLE_FORMATTED: Record<string, string> = {
  manager: 'Менеджер',
  cashier: 'Кассир',
  stocker: 'Кладовщик',
  // Add other roles as needed
}

// Status colors and formatted text
const STATUS_COLOR: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  disabled: 'bg-red-100 text-red-800',
}

const STATUS_FORMATTED: Record<string, string> = {
  active: 'Активен',
  disabled: 'Отключен',
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
