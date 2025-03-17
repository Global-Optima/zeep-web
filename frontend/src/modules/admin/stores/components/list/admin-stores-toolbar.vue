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
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
			>
				Экспорт
			</Button>
			<Button
				v-if="canCreate"
				@click="addStore"
			>
				Добавить
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { StoresFilter } from '@/modules/admin/stores/models/stores-dto.model'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter: StoresFilter }>()
const emit = defineEmits<{
  (event: 'update:filter', value: StoresFilter): void
}>()

// Local copy of the filter to avoid direct mutation
const localFilter = ref({ ...props.filter })

const canCreate = useHasRole(EmployeeRole.ADMIN)

// Search term with debouncing
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

// Watch debounced search term and emit updates with pagination reset
watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue
  emit('update:filter', {
    ...localFilter.value,
    search: newValue.trim(),
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
})

// Navigate to add store page
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_STORE_CREATE') })
}
</script>
