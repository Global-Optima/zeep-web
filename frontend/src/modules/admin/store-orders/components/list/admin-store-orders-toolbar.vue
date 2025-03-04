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

			<AdminSelectStoreDropdown
				v-if="showForFranchisee"
				:selected-store="selectedStore"
				@select="onSelectStore"
			/>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<AdminOrdersExport
				v-if="canExport"
				:store-id="filter.storeId"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Input } from '@/core/components/ui/input'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminOrdersExport from '@/modules/admin/store-orders/components/admin-orders-export.vue'
import type { OrdersFilterQuery } from '@/modules/admin/store-orders/models/orders.models'
import AdminSelectStoreDropdown from '@/modules/admin/stores/components/admin-select-store-dropdown.vue'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

// Props and Emit
const props = defineProps<{ filter: OrdersFilterQuery }>()
const emit = defineEmits<{(e: 'update:filter', value: OrdersFilterQuery): void }>()

const canExport = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])

const showForFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])

const selectedStore = ref<StoreDTO | undefined>(undefined)

const onSelectStore = (store: StoreDTO) => {
  selectedStore.value = store
  emit('update:filter', { ...props.filter, storeId: store.id})
}

// Local Filter
const localFilter = ref({ ...props.filter })

// Search Input
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

// Watch Search Input and Update Filter
watch(debouncedSearchTerm, (newValue) => {
	localFilter.value.search = newValue.trim()
	emit('update:filter', { ...localFilter.value })
})
</script>
