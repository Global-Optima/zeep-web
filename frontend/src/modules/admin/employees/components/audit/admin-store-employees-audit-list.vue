<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead>Активность</TableHead>
				<TableHead>Контекст</TableHead>
				<TableHead>Операция</TableHead>
				<TableHead>Дата</TableHead>
				<TableHead>Детали</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="(audit) in audits"
				:key="audit.id"
				class="h-12"
			>
				<!-- Display Activity Message -->
				<TableCell>
					<span v-html="formatLocalizedMessage(audit.localizedMessages.ru)"></span>
				</TableCell>

				<!-- Display Component Name -->
				<TableCell>
					{{ formattedComponentName(audit.componentName) }}
				</TableCell>

				<!-- Display Operation Type -->
				<TableCell>
					{{ formattedOperationType(audit.operationType) }}
				</TableCell>

				<!-- Display Date -->
				<TableCell>
					{{ formatDate(audit.timestamp) }}
				</TableCell>

				<!-- HoverCard for Details -->
				<TableCell>
					<HoverCard :open-delay="0">
						<HoverCardTrigger as-child>
							<Button
								variant="ghost"
								size="icon"
								class="p-0 w-9"
							>
								<SquareChartGantt class="text-gray-500 size-5" />
							</Button>
						</HoverCardTrigger>
						<HoverCardContent class="bg-gray-50 rounded-xl w-fit text-gray-800 text-sm">
							<pre>{{ formatDetails(audit.details) }}</pre>
						</HoverCardContent>
					</HoverCard>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/core/components/ui/hover-card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { formatLocalizedMessage } from '@/core/utils/format-localized-messages.utils'
import { EmployeeAuditComponentName, EmployeeAuditOperationType, FORMATTED_AUDIT_COMPONENTS, FORMATTED_AUDIT_OPERATION, type EmployeeAuditDTO } from '@/modules/admin/employees/models/employees-audit.models'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { SquareChartGantt } from 'lucide-vue-next'

defineProps<{
  audits: EmployeeAuditDTO[]
}>()

// Function to format the component name using the translations
const formattedComponentName = (componentName: EmployeeAuditComponentName): string => {
  return FORMATTED_AUDIT_COMPONENTS[componentName] || componentName
}

// Function to format the operation type using the translations
const formattedOperationType = (operationType: EmployeeAuditOperationType): string => {
  return FORMATTED_AUDIT_OPERATION[operationType] || operationType
}

// Function to format the date
const formatDate = (date: Date): string => {
  return format(new Date(date), 'dd MMM yyyy, HH:mm', { locale: ru })
}


// Function to filter and format the details object
const formatDetails = (details: Record<string, unknown>): string => {
  const filteredDetails = Object.fromEntries(
    Object.entries(details).filter(([_, value]) => value !== null && value !== undefined && value !== '')
  )
  return JSON.stringify(filteredDetails, null, 2)
}
</script>

<style scoped></style>
