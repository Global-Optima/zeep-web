<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Search Input -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Input
				v-model="searchTerm"
				placeholder="Поиск по названию"
				type="search"
				class="bg-white w-full md:w-64"
			/>
		</div>

		<!-- Status Filter -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Select
				v-model="selectedStatus"
				class="w-full md:w-64"
			>
				<SelectTrigger class="bg-white">
					<SelectValue placeholder="Фильтр по статусу" />
				</SelectTrigger>
				<SelectContent>
					<SelectItem value="all">Все</SelectItem>
					<SelectItem :value="ProjectStatus.PLANNED">Запланирован</SelectItem>
					<SelectItem :value="ProjectStatus.IN_PROGRESS">В процессе</SelectItem>
					<SelectItem :value="ProjectStatus.COMPLETED">Завершен</SelectItem>
				</SelectContent>
			</Select>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Input } from '@/core/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import { type ProjectFilter, ProjectStatus } from '@/modules/admin/admin-projects/models/projects.model'
import { ref, watch } from 'vue'

const props = defineProps<{ filter: ProjectFilter }>()

const emit = defineEmits(['update:filter'])

const localFilter = ref({ ...props.filter })

const searchTerm = ref(localFilter.value.search || '')
const selectedStatus = ref<ProjectStatus | 'all'>()

watch([searchTerm, selectedStatus], ([search, status]) => {
  emit('update:filter', { search: search.trim(), status: status === "all" ? null : status })
})
</script>
