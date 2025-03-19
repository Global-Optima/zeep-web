<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center gap-2 space-y-4 md:space-y-0 mb-4"
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
					<DropdownMenuItem @click="applyFilter('lowStock')">
						Только заканчивающиеся
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>

			<AdminSelectWarehouseDropdown
				v-if="showForRegion"
				:selected-warehouse="selectedWarehouse"
				@select="onSelectWarehouse"
			/>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
				>Экспорт</Button
			>
			<Button
				v-if="canCreate"
				@click="addStore"
				>Добавить</Button
			>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { GetWarehouseStockFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import AdminSelectWarehouseDropdown from '@/modules/admin/warehouses/components/admin-select-warehouse-dropdown.vue'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { useDebounce } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter: GetWarehouseStockFilter }>()
const emit = defineEmits<{ (e: 'update:filter', value: GetWarehouseStockFilter): void }>()

// Local Filter Copy
const localFilter = ref({ ...props.filter })

const showForRegion = useHasRole([EmployeeRole.REGION_WAREHOUSE_MANAGER])
const canCreate = useHasRole([EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER])

// Search Input with Debouncing
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

const selectedWarehouse = ref<WarehouseDTO | undefined>(undefined)

// Update filter when a warehouse is selected, with pagination reset
const onSelectWarehouse = (warehouse: WarehouseDTO) => {
  selectedWarehouse.value = warehouse
  emit('update:filter', {
    ...props.filter,
    warehouseId: warehouse.id,
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
}

// Watch debounced search input and update filter with pagination reset
watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue.trim()
  emit('update:filter', {
    ...localFilter.value,
    page: DEFAULT_PAGINATION_META.page,
    pageSize: DEFAULT_PAGINATION_META.pageSize
  })
})

// Filter Selection: update filter and reset pagination
const applyFilter = (filterType: string) => {
  if (filterType === 'all') {
    localFilter.value.lowStockOnly = undefined // Show all
  } else if (filterType === 'lowStock') {
    localFilter.value.lowStockOnly = true // Only low stock
  }
  emit('update:filter', {
    ...localFilter.value,
    page: DEFAULT_PAGINATION_META.page,
    pageSize: DEFAULT_PAGINATION_META.pageSize
  })
}

// Add Store Navigation
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_WAREHOUSE_STOCKS_CREATE') })
}
</script>
