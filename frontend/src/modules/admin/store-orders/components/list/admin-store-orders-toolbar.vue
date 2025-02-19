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
			<AdminOrdersExport v-if="canExport" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { Input } from '@/core/components/ui/input'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminOrdersExport from '@/modules/admin/store-orders/components/admin-orders-export.vue'
import type { OrdersFilterQuery } from '@/modules/admin/store-orders/models/orders.models'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

// Props and Emit
const props = defineProps<{ filter: OrdersFilterQuery }>()
const emit = defineEmits(['update:filter'])

const canExport = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])

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
