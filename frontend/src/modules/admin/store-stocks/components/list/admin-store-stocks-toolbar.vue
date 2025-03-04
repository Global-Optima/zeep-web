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

			<AdminSelectStoreDropdown
				v-if="showForFranchisee"
				:selected-store="selectedStore"
				@select="onSelectStore"
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
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { GetStoreWarehouseStockFilterQuery } from '@/modules/admin/store-stocks/models/store-stock.model'
import AdminSelectStoreDropdown from '@/modules/admin/stores/components/admin-select-store-dropdown.vue'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useDebounce } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

// Props and Emit
const props = defineProps<{ filter: GetStoreWarehouseStockFilterQuery }>()
const emit = defineEmits<{(e: 'update:filter', value: GetStoreWarehouseStockFilterQuery): void }>()

// Local Filter
const localFilter = ref({ ...props.filter })

// Search Input
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

const showForFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])
const canCreate = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA])

const selectedStore = ref<StoreDTO | undefined>(undefined)

const onSelectStore = (store: StoreDTO) => {
  selectedStore.value = store
  emit('update:filter', { ...props.filter, storeId: store.id})
}

// Watch Search Input and Update Filter
watch(debouncedSearchTerm, (newValue) => {
	localFilter.value.search = newValue.trim()
	emit('update:filter', { ...localFilter.value })
})

// Filter Selection
const applyFilter = (filterType: string) => {
	if (filterType === 'all') {
		localFilter.value.lowStockOnly = undefined // Show all
	} else if (filterType === 'lowStock') {
		localFilter.value.lowStockOnly = true // Only low stock
	}
	emit('update:filter', { ...localFilter.value })
}

// Add Store Navigation
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_CREATE_STORE_STOCKS') })
}
</script>
