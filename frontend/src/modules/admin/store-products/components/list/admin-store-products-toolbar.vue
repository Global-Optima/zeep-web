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
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { StoreProductsFilterDTO } from '@/modules/admin/store-products/models/store-products.model'
import AdminSelectStoreDropdown from '@/modules/admin/stores/components/admin-select-store-dropdown.vue'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ filter: StoreProductsFilterDTO }>()
const emit = defineEmits<{ (e: 'update:filter', value: StoreProductsFilterDTO): void }>()

// Local Filter Copy
const localFilter = ref({ ...props.filter })

// Search Input with Debouncing
const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

const showForFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])
const canCreate = useHasRole([EmployeeRole.STORE_MANAGER])
const selectedStore = ref<StoreDTO | undefined>(undefined)

// Update Filter when a store is selected, with pagination reset
const onSelectStore = (store: StoreDTO) => {
  selectedStore.value = store
  emit('update:filter', {
    ...props.filter,
    storeId: store.id,
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
}

// Watch debounced search input and update filter with pagination reset
watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue.trim()
  emit('update:filter', {
    ...localFilter.value,
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
})

// Add Store Navigation
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_STORE_PRODUCT_CREATE') })
}
</script>
