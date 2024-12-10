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
		<Table
			v-else
			class="bg-white rounded-xl"
		>
			<TableHeader>
				<TableRow>
					<TableHead class="p-4">Сотрудник</TableHead>
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
						<span class="font-medium"> {{ employee.firstName }} {{ employee.lastName }} </span>
					</TableCell>
					<!-- Role -->
					<TableCell class="p-4">
						{{ ROLE_FORMATTED[employee.role] }}
					</TableCell>

					<TableCell class="p-4">
						{{ formatPhoneNumber(employee.phone) }}
					</TableCell>

					<TableCell class="p-4">
						<span
							:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs ',
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
  status: string
}

const employees = ref<Employee[]>([
  {
    id: 1,
    firstName: 'Алексей',
    lastName: 'Смирнов',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'manager',
    phone: '+79161234567',
    status: 'active',
  },
  {
    id: 2,
    firstName: 'Екатерина',
    lastName: 'Иванова',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'cashier',
    phone: '+79169876543',
    status: 'disabled',
  },
  {
    id: 3,
    firstName: 'Дмитрий',
    lastName: 'Попов',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'barista',
    phone: '+79261234568',
    status: 'active',
  },
  {
    id: 4,
    firstName: 'Ольга',
    lastName: 'Кузнецова',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'cashier',
    phone: '+79269876544',
    status: 'active',
  },
  {
    id: 5,
    firstName: 'Игорь',
    lastName: 'Лебедев',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'manager',
    phone: '+79361234569',
    status: 'disabled',
  },
  {
    id: 6,
    firstName: 'Анна',
    lastName: 'Васильева',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'barista',
    phone: '+79369876545',
    status: 'active',
  },
  {
    id: 7,
    firstName: 'Максим',
    lastName: 'Соколов',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'technician',
    phone: '+79461234560',
    status: 'active',
  },
  {
    id: 8,
    firstName: 'Татьяна',
    lastName: 'Зайцева',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'cashier',
    phone: '+79469876546',
    status: 'disabled',
  },
  {
    id: 9,
    firstName: 'Сергей',
    lastName: 'Морозов',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'security',
    phone: '+79561234561',
    status: 'active',
  },
  {
    id: 10,
    firstName: 'Марина',
    lastName: 'Фёдорова',
    avatar: 'https://via.placeholder.com/32x32',
    role: 'manager',
    phone: '+79569876547',
    status: 'active',
  },
]);


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
  barista: 'Бариста',
  technician: 'Техник',
  security: 'Охранник',
};

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
