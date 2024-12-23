<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Input
				v-model="searchTerm"
				placeholder="Поиск"
				type="search"
				class="bg-white w-full md:w-64"
			/>
		</div>

		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button variant="outline">Экспорт</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import type { GetWarehouseStockFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

const props = defineProps<{ filter: GetWarehouseStockFilter }>()
const emit = defineEmits(['update:filter'])

const localFilter = ref({ ...props.filter })

const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

watch(debouncedSearchTerm, (newValue) => {
	localFilter.value.search = newValue.trim()
	emit('update:filter', { ...localFilter.value })
})
</script>
