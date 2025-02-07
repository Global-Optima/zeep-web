<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
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
				/>

				<!-- Material List -->
				<div class="max-h-[50vh] overflow-y-auto">
					<p
						v-if="!additives || additives.data.length === 0"
						class="text-muted-foreground"
					>
						Склады не найдены
					</p>

					<ul v-else>
						<li
							v-for="additive in additives.data"
							:key="additive.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(additive)"
						>
							<span>{{ additive.name }}</span>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="additives && additives.pagination.pageSize < additives.pagination.totalCount"
					variant="ghost"
					class="mt-4 w-full"
					@click="loadMore"
				>
					Еще
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

import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'

import type { WarehouseDTO, WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'

const {open} = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', additive: WarehouseDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<WarehouseFilter>({
  page: 1,
  pageSize: 10,
  search: ''
})


watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
  refetch()
})

const { data: additives, refetch } = useQuery({
  queryKey: computed(() => [
  'admin-warehouses',
  filter.value
]),
  queryFn: () => warehouseService.getPaginated(filter.value),
})


function loadMore() {
  if (!additives.value) return
  const pagination = additives.value.pagination

  if (pagination.pageSize < pagination.totalCount) {
    if(filter.value.pageSize) filter.value.pageSize += 10
  }
}

function selectMaterial(additive: WarehouseDTO) {
  emit('select', additive)
  onClose()
}

function onClose() {
  filter.value = {
    page: 1,
    pageSize: 10,
    search: ''
  }

  emit('close')
}
</script>

<style scoped></style>
