<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название проекта</TableHead>
				<TableHead>Дата начала</TableHead>
				<TableHead>Дата завершения</TableHead>
				<TableHead>Статус</TableHead>
				<TableHead>Ответственный</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="project in projects"
				:key="project.id"
				@click="onProjectClick(project.id)"
				class="hover:bg-slate-50 transition-all duration-200 cursor-pointer"
				:title="`Открыть детали проекта: ${project.name}`"
			>
				<TableCell class="py-4 font-medium">{{ project.name }}</TableCell>
				<TableCell>{{ project.startDate }}</TableCell>
				<TableCell>{{ project.endDate }}</TableCell>
				<TableCell>
					<p
						class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
						:class="getStatusClass(project.status)"
					>
						{{project.status}}
					</p>
				</TableCell>
				<TableCell>
					{{ project.responsible }}
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
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
import { type Project, ProjectStatus } from '@/modules/admin/admin-projects/models/projects.model'
import { useRouter } from 'vue-router'

const { projects } = defineProps<{ projects: Project[] }>()

const router = useRouter()

const STATUS_COLORS: Record<ProjectStatus, string> = {
  [ProjectStatus.PLANNED]: 'bg-yellow-100 text-yellow-800',
  [ProjectStatus.IN_PROGRESS]: 'bg-blue-100 text-blue-800',
  [ProjectStatus.COMPLETED]: 'bg-green-100 text-green-800'
}

const getStatusClass = (status: ProjectStatus): string => {
  console.log(status)
	return STATUS_COLORS[status] || ''
}

const onProjectClick = (id: number) => {
	router.push(`/admin/projects/${id}`)
}
</script>

<style scoped></style>
