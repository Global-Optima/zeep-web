<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Search Input and Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<!-- Search Input -->
			<Input
				v-model="searchTerm"
				placeholder="Поиск"
				type="search"
				class="bg-white w-full md:w-64"
			/>

			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button
						variant="outline"
						class="whitespace-nowrap"
					>
						Фильтр
						<ChevronDown class="ml-2 w-4 h-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem @click="applyFilter('all')"> Все </DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter('active')"> Активные </DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter('disabled')"> Деактивированные </DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
				>Экспорт</Button
			>
			<Button @click="onCreateClick"> Создать </Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import { useDebounce } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter: EmployeesFilter }>()
const emit = defineEmits<{
  (event: 'update:filter', value: EmployeesFilter): void
}>()

const router = useRouter()

// Local Filter: Create a local copy of the filter to avoid direct mutation
const localFilter = ref({ ...props.filter })

// Search Input: Initialize the search term with any existing value
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

// Watch Search Input and Update Filter
watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue.trim()
  emit('update:filter', {
    ...localFilter.value,
    search: newValue.trim(),
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
})

// Filter Selection
const applyFilter = (filterType: string) => {
	if (filterType === 'all') {
		localFilter.value.isActive = undefined
	} else if (filterType === 'active') {
		localFilter.value.isActive = true
	}
  else if (filterType === 'disabled') {
		localFilter.value.isActive = false
	}
	emit('update:filter', { ...localFilter.value })
}

const onCreateClick = () => {
  router.push({name: getRouteName("ADMIN_ADMIN_EMPLOYEE_CREATE")})
}
</script>
