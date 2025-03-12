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
							@click="$router.push(`/admin/employees/${employee.employeeId}/audit`)"
							>Еще</Button
						>
					</div>
				</CardHeader>
				<CardContent>
					<AdminEmployeesDetailsAudits :audits="employeeAudits?.data ?? []" />
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
import AdminEmployeesDetailsAudits from '@/modules/admin/employees/components/details/admin-employees-details-audits.vue'
import AdminEmployeesDetailsInfo from '@/modules/admin/employees/components/details/admin-employees-details-info.vue'
import AdminEmployeesDetailsShifts from '@/modules/admin/employees/components/details/admin-employees-details-shifts.vue'
import { employeeAuditService } from '@/modules/admin/employees/services/employees-audit.service'
import { employeesService } from '@/modules/admin/employees/services/employees.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

interface Shift {
  date: string;
  shift: string;
}

const { data: employee } = useQuery({
	queryKey: ['employee-current'],
	queryFn: () => employeesService.getCurrentEmployee(),
})

const { data: employeeAudits } = useQuery({
	queryKey: computed(() => ['employee-audits', employee.value?.id]),
	queryFn: () => employeeAuditService.getAudits({ employeeId: Number(employee.value?.id) }),
	enabled: computed(() => Boolean(employee.value?.id)),
})

const employeeShifts = ref<Shift[]>([
  { date: 'Понедельник', shift: '9:00 - 17:00' },
  { date: 'Вторник', shift: '10:00 - 18:00' },
  { date: 'Среда', shift: '9:00 - 17:00' },
  { date: 'Четверг', shift: '11:00 - 19:00' },
  { date: 'Пятница', shift: '9:00 - 17:00' },
  { date: 'Суббота', shift: 'Выходной' },
  { date: 'Воскресенье', shift: 'Выходной' },
])
</script>

<style scoped></style>
