<template>
	<AdminProjectsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="filteredProjects.length === 0"
				class="text-muted-foreground"
			>
				Проекты не найдены
			</p>
			<AdminProjectsList
				v-else
				:projects="filteredProjects"
			/>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminProjectsList from '@/modules/admin/admin-projects/components/list/admin-projects-list.vue'
import AdminProjectsToolbar from '@/modules/admin/admin-projects/components/list/admin-projects-toolbar.vue'
import { ProjectStatus, type Project, type ProjectFilter } from '@/modules/admin/admin-projects/models/projects.model'
import { computed, ref } from 'vue'

const filter = ref<ProjectFilter>({
  search: '',
  status: null, // Status filter
})

// Mock Data for Projects
const projects = ref<Project[]>([
  {
    id: 1,
    name: 'Проект ремонта оборудования',
    startDate: '2024-01-01',
    endDate: '2024-06-01',
    status: ProjectStatus.IN_PROGRESS,
    responsible: 'Иван Иванов',
  },
  {
    id: 2,
    name: 'Модернизация системы учета',
    startDate: '2023-05-01',
    endDate: '2023-12-31',
    status: ProjectStatus.COMPLETED,
    responsible: 'Анна Смирнова',
  },
  {
    id: 3,
    name: 'Поставка оборудования для склада',
    startDate: '2024-02-15',
    endDate: '2024-05-15',
    status: ProjectStatus.PLANNED,
    responsible: 'Сергей Павлов',
  },
  {
    id: 4,
    name: 'Обучение персонала',
    startDate: '2024-03-01',
    endDate: '2024-03-30',
    status: ProjectStatus.IN_PROGRESS,
    responsible: 'Мария Федорова',
  },
])

const filteredProjects = computed(() => {
  return projects.value.filter((project) => {
    const matchesSearch =
      !filter.value.search ||
      project.name.toLowerCase().includes(filter.value.search.toLowerCase())
    const matchesStatus =
      !filter.value.status || project.status === filter.value.status
    return matchesSearch && matchesStatus
  })
})

function updateFilter(updatedFilter: ProjectFilter) {
  filter.value = { ...filter.value, ...updatedFilter }
}
</script>

<style scoped></style>
