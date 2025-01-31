<template>
	<div>
		<!-- Display Message if No Activities -->
		<p
			v-if="audits.length === 0"
			class="py-4 text-center text-muted-foreground"
		>
			Активности отсутствуют
		</p>

		<!-- Table of Activities -->
		<Table
			v-else
			class="bg-white rounded-xl"
		>
			<TableHeader>
				<TableRow>
					<TableHead>Активность</TableHead>
					<TableHead>Дата</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<TableRow
					v-for="(audit, index) in audits"
					:key="index"
					class="h-12"
				>
					<TableCell>
						<span v-html="formatLocalizedMessage(audit.localizedMessages.ru)"></span>
					</TableCell>

					<TableCell>
						{{ formatDate(audit.timestamp) }}
					</TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
</template>

<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { formatLocalizedMessage } from '@/core/utils/format-localized-messages.utils'
import { type EmployeeAuditDTO } from '@/modules/admin/store-employees/models/employees-audit.models'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
defineProps<{
  audits: EmployeeAuditDTO[]
}>()

const formatDate = (date: Date): string => {
  return format(new Date(date), 'dd MMM yyyy, HH:mm', { locale: ru })
}
</script>

<style scoped></style>
