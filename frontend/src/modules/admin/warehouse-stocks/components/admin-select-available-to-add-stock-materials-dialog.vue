<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите доступный материал</DialogTitle>
			</DialogHeader>

			<div>
				<Input
					v-model="searchTerm"
					placeholder="Поиск"
					type="search"
					class="mt-2 mb-4 w-full"
				/>

				<div class="max-h-[50vh] overflow-y-auto">
					<p
						v-if="!stockMaterials || stockMaterials.data.length === 0"
						class="text-muted-foreground"
					>
						Материалы не найдены
					</p>

					<ul v-else>
						<li
							v-for="stockMaterial in stockMaterials.data"
							:key="stockMaterial.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(stockMaterial)"
						>
							<span class="flex-1">{{ stockMaterial.name }}</span>
							<span class="text-gray-500 text-sm">
								{{ stockMaterial.category.name }}
							</span>
						</li>
					</ul>
				</div>

				<Button
					v-if="showLoadMoreButton"
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
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { AvailableWarehouseStockMaterialsFilter } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, reactive, ref } from 'vue'

const { open } = defineProps<{
	open: boolean
}>()

const emit = defineEmits<{
	(e: 'close'): void
	(e: 'select', stockMaterial: StockMaterialsDTO): void
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(searchTerm, 500)
const localFilter = reactive<Partial<AvailableWarehouseStockMaterialsFilter>>({
	page: 1,
	pageSize: 10,
})
const mergedFilter = computed<AvailableWarehouseStockMaterialsFilter>(() => ({
	...localFilter,
	search: debouncedSearchTerm.value.trim(),
}))

const { data: stockMaterials } = useQuery({
	queryKey: computed(() => ['available-stock-materials', mergedFilter.value]),
	queryFn: () => warehouseStocksService.getAvailableStockMaterials(mergedFilter.value),
	enabled: computed(() => open),
})

function loadMore() {
	if (!stockMaterials.value) return
	const pagination = stockMaterials.value.pagination
	if (pagination.pageSize < pagination.totalCount) {
		localFilter.pageSize = (localFilter.pageSize ?? 10) + 10
	}
}

function selectMaterial(stockMaterial: StockMaterialsDTO) {
	emit('select', stockMaterial)
	onClose()
}

function onClose() {
	searchTerm.value = ''
	localFilter.page = 1
	localFilter.pageSize = 10
	emit('close')
}

const showLoadMoreButton = computed(() => {
	if (!stockMaterials.value) return false
	const { pageSize, totalCount } = stockMaterials.value.pagination
	return pageSize < totalCount
})
</script>

<style scoped></style>
