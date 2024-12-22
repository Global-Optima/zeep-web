<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите материал</DialogTitle>
			</DialogHeader>
			<DialogBody>
				<!-- Search Input -->
				<Input
					v-model="searchTerm"
					placeholder="Поиск материала"
					type="search"
					class="mt-2 mb-4 w-full"
				/>

				<!-- Material List -->
				<div class="max-h-64 overflow-y-auto">
					<p
						v-if="materials.length === 0"
						class="text-muted-foreground"
					>
						Товары не найдены
					</p>

					<ul v-else>
						<li
							v-for="material in materials"
							:key="material.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(material)"
						>
							<span>{{ material.name }}</span>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="pagination.page < pagination.totalPages"
					variant="outline"
					class="mt-4 w-full"
					@click="loadMore"
				>
					Еще
				</Button>
			</DialogBody>
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

import type { PaginationMeta } from '@/core/utils/pagination.utils'
import type { StockMaterialsDTO, StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'

const {open} = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', material: StockMaterialsDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<StockMaterialsFilter>({
  page: 1,
  pageSize: 10,
  search: ''
})

const materials = ref<StockMaterialsDTO[]>([])

const pagination = ref<PaginationMeta>({
  page: 0,
  pageSize: 0,
  totalCount: 0,
  totalPages: 0
})


watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
  refetch()
})


const { data: queryData, refetch } = useQuery({
  queryKey: computed(() => [
  'stock-materials',
  filter.value
]),
  queryFn: () => stockMaterialsService.getAllStockMaterials(filter.value),
})


watch(queryData, (newData) => {
  if (!newData) return

  pagination.value = newData.pagination

  if (pagination.value.page === 1) {
    materials.value = newData.data
  } else {
    materials.value = materials.value.concat(newData.data)
  }
})

function loadMore() {
  if (pagination.value.page < pagination.value.totalPages) {
    if(filter.value.page) filter.value.page += 1
    refetch()
  }
}

function selectMaterial(material: StockMaterialsDTO) {
  emit('select', material)
  onClose()
}

function onClose() {
  filter.value = {
    page: 1,
    pageSize: 10,
    search: ''
  }

  pagination.value = {
    page: 0,
    pageSize: 0,
    totalCount: 0,
    totalPages: 0
  }

  emit('close')
}
</script>

<style scoped></style>
