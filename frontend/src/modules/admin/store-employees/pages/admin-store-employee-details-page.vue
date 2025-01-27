<template>
	<p
		v-if="!employee"
		class="text-muted-foreground"
	>
		Сотрудник не найден
	</p>

	<div
		v-else
		class="flex md:flex-row flex-col gap-4 mx-auto max-w-7xl"
	>
		<div class="w-full md:w-1/3">
			<AdminEmployeesDetailsInfo :employee="employee" />
		</div>

		<div class="flex flex-col gap-4 w-full md:w-2/3">
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

			<Card>
				<CardHeader>
					<CardTitle> График рабочих смен</CardTitle>
					<CardDescription>Просмотр смен сотрудника</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminEmployeesDetailsShifts :shifts="employeeShifts" />
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<div class="flex justify-between items-start gap-4">
						<div>
							<CardTitle>Недавние активности</CardTitle>
							<CardDescription>Последние выполненные задачи сотрудника</CardDescription>
						</div>

						<Button
							variant="outline"
							@click="$router.push(`/admin/store-employees/${employeeId}/audit`)"
							>Еще</Button
						>
					</div>
				</CardHeader>
				<CardContent>
					<AdminEmployeesDetailsActivities :audits="employeeAudits?.data ?? []" />
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { formatPrice } from '@/core/utils/price.utils'
import AdminEmployeesDetailsActivities from '@/modules/admin/store-employees/components/details/admin-employees-details-activities.vue'
import AdminEmployeesDetailsInfo from '@/modules/admin/store-employees/components/details/admin-employees-details-info.vue'
import AdminEmployeesDetailsShifts from '@/modules/admin/store-employees/components/details/admin-employees-details-shifts.vue'
import AdminEmployeesDetailsStats from '@/modules/admin/store-employees/components/details/admin-employees-details-stats.vue'
import { employeeAuditService } from '@/modules/admin/store-employees/services/employees-audit.service'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRoute } from 'vue-router'

interface Stat {
  totalSales: string;
  hoursWorked: string;
  tasksCompleted: string;
}

interface Shift {
  date: string;
  shift: string;
}

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
  queryKey: ['store-employee', employeeId],
	queryFn: () => employeesService.getStoreEmployeeById(Number(employeeId)),
  enabled: !!employeeId,
})


const { data: employeeAudits } = useQuery({
  queryKey: ['employee-audits', employeeId],
	queryFn: () => employeeAuditService.getAudits({employeeId: Number(employeeId)}),
  enabled: !!employeeId,
})

const employeeStats = ref<Stat>({
  totalSales: formatPrice(3500000),
  hoursWorked: '1 200 часов',
  tasksCompleted: '350 задач',
});


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

<style scoped></style>
