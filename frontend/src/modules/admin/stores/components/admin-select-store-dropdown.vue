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
				<DialogTitle>Выберите кафе</DialogTitle>
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
						v-if="!stores || stores.data.length === 0"
						class="py-4 text-muted-foreground text-center"
					>
						Кафе не найдены
					</p>

					<ul v-else>
						<li
							v-for="warehouse in stores.data"
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
import { computed, ref, watch } from 'vue'

// Components
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { ChevronDown } from 'lucide-vue-next'

// Types and Services
import type { StoresFilter } from '@/modules/admin/stores/models/stores-dto.model'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { storesService } from '@/modules/admin/stores/services/stores.service'

// Props and Emits
const {selectedStore,initialFilter, selectFirst = true} = defineProps<{
  selectedStore?: StoreDTO
  initialFilter?: StoresFilter
  selectFirst?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', store: StoreDTO): void
}>()

// State Management
const isOpen = ref(false)
const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(searchTerm, 500)
const filter = ref<StoresFilter>({
  page: 1,
  pageSize: 10,
  search: initialFilter?.search ?? ''
})

// Query for Warehouses
const { data: stores } = useQuery({
  queryKey: computed(() => [
  'admin-stores',
  filter.value
]),
  queryFn: () => storesService.getPaginated(filter.value),
})

// Computed Properties
const canLoadMore = computed(() => {
  if (!stores.value) return false
  const pagination = stores.value.pagination
  return pagination.pageSize < pagination.totalCount
})

const buttonLabel = computed(() => selectedStore?.name ?? "Кафе не выбран")

const selectedWarehouseId = computed(() => selectedStore?.id)

// Watchers
watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
})


watch(stores, (newStores) => {
  if (selectFirst && newStores?.data.length) {
    selectWarehouse(newStores.data[0])
  }
}, { immediate: true })

// Methods
function loadMore() {
  if (!stores.value) return
  if (filter.value.page) {
    filter.value = {
    ...filter.value,
    page: filter.value.page + 1
  }
  }
}

function selectWarehouse(warehouse: StoreDTO) {
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
