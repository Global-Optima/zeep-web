<template>
	<div class="space-y-6 mx-auto w-full max-w-6xl">
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				type="button"
				@click="$router.back"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>

			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Детали проекта
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="$router.back"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					disabled
				>
					Сохранить
				</Button>
			</div>
		</div>

		<Card>
			<CardHeader>
				<CardTitle class="text-lg sm:text-xl">{{ project.name }}</CardTitle>
				<CardDescription> {{project.description}} </CardDescription>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Дата начала</TableHead>
							<TableHead>Дата завершения</TableHead>
							<TableHead>Ответственный</TableHead>
							<TableHead>Статус</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow class="hover:bg-slate-50 transition-all duration-200 cursor-pointer">
							<TableCell>{{ project.startDate }}</TableCell>
							<TableCell>{{ project.endDate }}</TableCell>
							<TableCell>
								{{ project.responsible }}
							</TableCell>
							<TableCell>
								<p
									class="inline-flex items-center rounded-md w-fit"
									:class="getStatusClass(project.status)"
								>
									{{project.status}}
								</p>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Products Table -->
		<Card>
			<CardHeader>
				<CardTitle class="text-lg sm:text-xl">Активные товары в проекте</CardTitle>
				<CardDescription>Список товаров, задействованных в этом проекте</CardDescription>
			</CardHeader>
			<CardContent>
				<AdminProductsPage />
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { ProjectStatus, type Project } from '@/modules/admin/admin-projects/models/projects.model'
import AdminProductsPage from '@/modules/admin/products/pages/admin-products-page.vue'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// Mock project details data
const mockProjectDetails: Record<number, Project & { description: string }> = {
	1: {
		id: 1,
		name: 'Проект ремонта оборудования',
		startDate: '2024-01-01',
		endDate: '2024-06-01',
		status: ProjectStatus.IN_PROGRESS,
		responsible: 'Иван Иванов',
		description: 'Детальное описание проекта ремонта оборудования, включая ключевые задачи и сроки.',
	},
	2: {
		id: 2,
		name: 'Модернизация системы учета',
		startDate: '2023-05-01',
		endDate: '2023-12-31',
		status: ProjectStatus.COMPLETED,
		responsible: 'Анна Смирнова',
		description: 'Проект по модернизации системы учета, направленный на улучшение процесса управления данными.',
	},
}

const projectId = Number(route.params.id)
const project = ref(mockProjectDetails[projectId])

if (!project.value) {
	router.replace('/admin/projects') // Redirect if project doesn't exist
}

const STATUS_COLORS: Record<ProjectStatus, string> = {
	[ProjectStatus.PLANNED]: 'text-yellow-700',
	[ProjectStatus.IN_PROGRESS]: 'text-blue-700',
	[ProjectStatus.COMPLETED]: 'text-green-700',
}


const getStatusClass = (status: ProjectStatus): string => STATUS_COLORS[status] || ''
</script>
