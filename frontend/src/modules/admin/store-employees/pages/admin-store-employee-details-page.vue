<template>
	<p
		v-if="!employee"
		class="text-muted-foreground"
	>
		Сотрудник не найден
	</p>

	<div
		v-else
		class="flex md:flex-row flex-col gap-4"
	>
		<!-- Left Side: Employee Details Card -->
		<div class="w-full md:w-1/3">
			<AdminEmployeesDetailsInfo :employee="employee" />
		</div>

		<!-- Right Side: Main Content -->
		<div class="flex flex-col gap-4 w-full md:w-2/3">
			<!-- Statistics Cards -->
			<div class="gap-6 grid grid-cols-1 sm:grid-cols-3">
				<AdminEmployeesDetailsStats
					title="Общие продажи"
					:value="employeeStats.totalSales"
					icon="dollar-sign"
				/>
				<AdminEmployeesDetailsStats
					title="Отработанные часы"
					:value="employeeStats.hoursWorked"
					icon="clock"
				/>
				<AdminEmployeesDetailsStats
					title="Завершенные задачи"
					:value="employeeStats.tasksCompleted"
					icon="check-circle"
				/>
			</div>

			<!-- Activities Table -->
			<Card>
				<CardHeader>
					<CardTitle class="font-medium text-lg"> Недавние активности</CardTitle>
					<CardDescription>Последние выполненные задачи сотрудника</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminEmployeesDetailsActivities :activities="employeeActivities" />
				</CardContent>
			</Card>

			<!-- Working Shifts Calendar Card -->
			<Card>
				<CardHeader>
					<CardTitle class="font-medium text-lg"> График рабочих смен</CardTitle>
					<CardDescription>Просмотр смен сотрудника</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminEmployeesDetailsShifts :shifts="employeeShifts" />
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { formatPrice } from '@/core/utils/price.utils'
import AdminEmployeesDetailsActivities from '@/modules/admin/employees/components/details/admin-employees-details-activities.vue'
import AdminEmployeesDetailsInfo from '@/modules/admin/employees/components/details/admin-employees-details-info.vue'
import AdminEmployeesDetailsShifts from '@/modules/admin/employees/components/details/admin-employees-details-shifts.vue'
import AdminEmployeesDetailsStats from '@/modules/admin/employees/components/details/admin-employees-details-stats.vue'
import { employeesService } from '@/modules/admin/employees/services/employees.service'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRoute } from 'vue-router'

interface Stat {
  totalSales: string;
  hoursWorked: string;
  tasksCompleted: string;
}

interface Activity {
  date: string;
  activity: string;
  status: string;
}

interface Shift {
  date: string;
  shift: string;
}

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
  queryKey: ['employee', employeeId],
	queryFn: () => employeesService.getStoreEmployeeById(Number(employeeId)),
  enabled: !!employeeId,
})

const employeeStats = ref<Stat>({
  totalSales: formatPrice(3500000),
  hoursWorked: '1 200 часов',
  tasksCompleted: '350 задач',
});

const employeeActivities = ref<Activity[]>([
  { date: '2023-10-01', activity: 'Заключил сделку с компанией XYZ', status: 'Завершено' },
  { date: '2023-09-28', activity: 'Участвовал в собрании отдела продаж', status: 'Завершено' },
  { date: '2023-09-25', activity: 'Обновил базу данных клиентов', status: 'Завершено' },
  { date: '2023-09-22', activity: 'Отправил еженедельный отчет руководству', status: 'Завершено' },
  { date: '2023-09-20', activity: 'Провел переговоры с новым клиентом', status: 'Завершено' },
  { date: '2023-09-18', activity: 'Решил проблему клиента с заказом', status: 'Завершено' },
  { date: '2023-09-15', activity: 'Участвовал в тренинге по продажам', status: 'Завершено' },
  { date: '2023-09-13', activity: 'Подготовил презентацию для нового продукта', status: 'Завершено' },
]);

const employeeShifts = ref<Shift[]>([
  { date: 'Понедельник', shift: '9:00 - 17:00' },
  { date: 'Вторник', shift: '10:00 - 18:00' },
  { date: 'Среда', shift: '9:00 - 17:00' },
  { date: 'Четверг', shift: '11:00 - 19:00' },
  { date: 'Пятница', shift: '9:00 - 17:00' },
  { date: 'Суббота', shift: 'Выходной' },
  { date: 'Воскресенье', shift: 'Выходной' },
]);
</script>

<style scoped>
/* Add any additional styles if needed */
</style>
