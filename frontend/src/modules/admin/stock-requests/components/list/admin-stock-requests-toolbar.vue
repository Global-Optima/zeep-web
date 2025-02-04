<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Input
				v-model="searchTerm"
				placeholder="Поиск"
				type="search"
				class="bg-white w-full md:w-64"
			/>

			<MultiSelectFilter
				title="Статусы"
				:options="filteredStatusOptions"
				v-model="selectedStatuses"
			/>
		</div>

		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
			>
				Экспорт
			</Button>
			<Button @click="onCreateClick">Создать</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import MultiSelectFilter from '@/core/components/multi-select-filter/MultiSelectFilter.vue'
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import { STOCK_REQUEST_STATUS_OPTIONS, StockRequestStatus, type GetStockRequestsFilter } from '@/modules/admin/stock-requests/models/stock-requests.model'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter?: GetStockRequestsFilter }>()
const emit = defineEmits(['update:filter'])
const router = useRouter()

const localFilter = ref({ ...props.filter })

const selectedStatuses = ref<StockRequestStatus[]>(props.filter?.statuses ?? [])
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

const { currentEmployee } = useEmployeeAuthStore()

const filteredStatusOptions = computed(() => {
  if (!currentEmployee || !currentEmployee.role) return []

  const warehouseRoles: EmployeeRole[] = [EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER]
  const storeRoles: EmployeeRole[] = [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA]

  if (warehouseRoles.includes(currentEmployee.role)) {
    return STOCK_REQUEST_STATUS_OPTIONS.filter(status => status.value !== StockRequestStatus.CREATED)
  }

  if (storeRoles.includes(currentEmployee.role)) {
    return STOCK_REQUEST_STATUS_OPTIONS
  }

  return []
})

const filteredSelectedStatuses = computed(() => {
  return selectedStatuses.value.filter(status =>
    filteredStatusOptions.value.some(option => option.value === status)
  )
})

watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue
  emit('update:filter', { ...props.filter, search: newValue.trim() })
})

watch(filteredSelectedStatuses, (newStatuses) => {
  emit('update:filter', {
    ...props.filter,
    statuses: newStatuses.length ? newStatuses : undefined
  })
})

const onCreateClick = () => {
  router.push({ name: getRouteName('ADMIN_STORE_STOCK_REQUESTS_CREATE') })
}
</script>
