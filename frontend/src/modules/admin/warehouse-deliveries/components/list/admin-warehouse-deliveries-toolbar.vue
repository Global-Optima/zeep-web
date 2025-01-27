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
				>Экспорт</Button
			>
			<Button @click="addStore">Создать</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import type { WarehouseDeliveryFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

// Props and Emit
const props = defineProps<{ filter: WarehouseDeliveryFilter }>()
const emit = defineEmits(['update:filter'])

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

// Add Store Navigation
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_WAREHOUSE_DELIVERIES_CREATE') })
}
</script>
