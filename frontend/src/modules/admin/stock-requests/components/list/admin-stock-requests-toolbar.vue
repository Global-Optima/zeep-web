<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center gap-2 space-y-4 md:space-y-0 mb-4"
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

			<AdminSelectWarehouseDropdown
				v-if="showForRegion"
				:selected-warehouse="selectedWarehouse"
				@select="onSelectWarehouse"
			/>

			<AdminSelectStoreDropdown
				v-if="showForFranchisee"
				:selected-store="selectedStore"
				@select="onSelectStore"
			/>
		</div>

		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
			>
				Экспорт
			</Button>
			<Button
				v-if="canCreateStockRequest"
				@click="onCreateClick"
				>Создать</Button
			>
		</div>
	</div>
</template>

<script setup lang="ts">
import MultiSelectFilter from '@/core/components/multi-select-filter/MultiSelectFilter.vue'
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import { STOCK_REQUEST_STATUS_OPTIONS, StockRequestStatus, type GetStockRequestsFilter } from '@/modules/admin/stock-requests/models/stock-requests.model'
import AdminSelectStoreDropdown from '@/modules/admin/stores/components/admin-select-store-dropdown.vue'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import AdminSelectWarehouseDropdown from '@/modules/admin/warehouses/components/admin-select-warehouse-dropdown.vue'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter?: GetStockRequestsFilter }>()
const emit = defineEmits<{(e: 'update:filter', value: GetStockRequestsFilter): void }>()
const router = useRouter()

const canCreateStockRequest = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA])

const localFilter = ref({ ...props.filter })

const showForRegion = useHasRole([EmployeeRole.REGION_WAREHOUSE_MANAGER])

const selectedStatuses = ref<StockRequestStatus[]>(props.filter?.statuses ?? [])
const selectedWarehouse = ref<WarehouseDTO | undefined>(undefined)
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

const showForFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])

const selectedStore = ref<StoreDTO | undefined>(undefined)

const onSelectStore = (store: StoreDTO) => {
  selectedStore.value = store
  emit('update:filter', { ...props.filter, storeId: store.id})
}

const { currentEmployee } = useEmployeeAuthStore()

const onSelectWarehouse = (warehouse: WarehouseDTO) => {
  selectedWarehouse.value = warehouse
  emit('update:filter', { ...props.filter, warehouseId: warehouse.id})
}

const filteredStatusOptions = computed(() => {
  if (!currentEmployee || !currentEmployee.role) return []

  const warehouseRoles: EmployeeRole[] = [EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.REGION_WAREHOUSE_MANAGER]
  const storeRoles: EmployeeRole[] = [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA, EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER]

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
  if (!canCreateStockRequest) return
  router.push({ name: getRouteName('ADMIN_STORE_STOCK_REQUESTS_CREATE') })
}
</script>
