<template>
	<Dialog
		v-model:open="isOpen"
		@update:open="v => isOpen = v"
	>
		<DialogTrigger as-child>
			<Button
				variant="outline"
				class="justify-between gap-4 w-full"
				@click="isOpen = true"
			>
				<p>{{ buttonLabel }}</p>
				<ChevronDown class="size-4 text-gray-700" />
			</Button>
		</DialogTrigger>

		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите склад</DialogTitle>
			</DialogHeader>

			<div>
				<!-- Search Input -->
				<Input
					v-model="searchTerm"
					placeholder="Поиск"
					type="search"
					class="mt-2 mb-4 w-full"
					autofocus
				/>

				<!-- Warehouse List -->
				<div class="max-h-[50vh] overflow-y-auto">
					<p
						v-if="!warehouses || warehouses.data.length === 0"
						class="py-4 text-muted-foreground text-center"
					>
						Склады не найдены
					</p>

					<ul v-else>
						<li
							v-for="warehouse in warehouses.data"
							:key="warehouse.id"
							:class="[
                'flex items-center px-3 py-2 rounded-lg cursor-pointer',
                warehouse.id === selectedWarehouseId ? 'bg-green-50 dark:bg-gray-800' : 'hover:bg-gray-100 dark:hover:bg-gray-700'
              ]"
							@click="selectWarehouse(warehouse)"
							@keydown.enter="selectWarehouse(warehouse)"
							tabindex="0"
						>
							<span>{{ warehouse.name }}</span>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="canLoadMore"
					variant="ghost"
					class="mt-4 w-full"
					@click="loadMore"
				>
					Загрузить еще
				</Button>
			</div>

			<DialogFooter>
				<Button
					variant="outline"
					@click="onClose"
				>
					Закрыть
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, onMounted, ref, watch } from 'vue'

// Components
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { ChevronDown } from 'lucide-vue-next'

// Types and Services
import type { WarehouseDTO, WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'

// Props and Emits
const {selectedWarehouse,initialFilter, selectFirst = true} = defineProps<{
  selectedWarehouse?: WarehouseDTO
  initialFilter?: WarehouseFilter
  selectFirst?: boolean
}>()

const emit = defineEmits<{
(e: 'select', warehouse: WarehouseDTO): void
(e: 'update:open', value: boolean): void
}>()

// State Management
const isOpen = ref(false)
const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(searchTerm, 500)
const filter = ref<WarehouseFilter>({
  page: 1,
  pageSize: 10,
  search: initialFilter?.search ?? ''
})

// Query for Warehouses
const { data: warehouses, refetch } = useQuery({
queryKey: computed(() => ['admin-warehouses', filter.value]),
queryFn: () => warehouseService.getPaginated(filter.value),
})

// Computed Properties
const canLoadMore = computed(() => {
  if (!warehouses.value) return false
  const pagination = warehouses.value.pagination
  return pagination.pageSize < pagination.totalCount
})

const buttonLabel = computed(() => selectedWarehouse?.name ?? "Склад не выбран")

const selectedWarehouseId = computed(() => selectedWarehouse?.id)

// Watchers
watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
  refetch()
})


onMounted(() => {
  if (selectFirst  && warehouses.value?.data.length) {
    selectWarehouse(warehouses.value.data[0])
  }
})

// Methods
function loadMore() {
  if (!warehouses.value) return
  if (filter.value.page) {
    filter.value = {
    ...filter.value,
    page: filter.value.page + 1
  }
  }
}

function selectWarehouse(warehouse: WarehouseDTO) {
emit('select', warehouse)
isOpen.value = false
}

function onClose() {
  filter.value = {
    page: 1,
    pageSize: 10,
    search: initialFilter?.search ?? ''
  }
  searchTerm.value = ''
  isOpen.value = false
  }
</script>

<style scoped>
/* Add any scoped styles here */
</style>
