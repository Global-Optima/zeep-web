<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите упаковку</DialogTitle>
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
						v-if="!stockMaterialPackages || stockMaterialPackages.data.length === 0"
						class="text-muted-foreground"
					>
						Упаковки не найдены
					</p>

					<ul v-else>
						<li
							v-for="stockMaterial in stockMaterialPackages.data"
							:key="stockMaterial.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(stockMaterial)"
						>
							<span class="flex-1">{{ stockMaterial.size }} {{ stockMaterial.unit.name }}</span>
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
import type { StockMaterialPackageFilterDTO, StockMaterialPackagesDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, reactive, ref } from 'vue'

const { open, initialFilter } = defineProps<{
	open: boolean
	initialFilter?: StockMaterialPackageFilterDTO
}>()

const emit = defineEmits<{
	(e: 'close'): void
	(e: 'select', materialPackage: StockMaterialPackagesDTO): void
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(searchTerm, 500)
const localFilter = reactive<Partial<StockMaterialPackageFilterDTO>>({})
const mergedFilter = computed<StockMaterialPackageFilterDTO>(() => ({
	...(initialFilter ?? {}),
	...localFilter,
	search: debouncedSearchTerm.value.trim(),
}))

const { data: stockMaterialPackages } = useQuery({
	queryKey: computed(() => ['stock-material-packages', mergedFilter.value]),
	queryFn: () => stockMaterialsService.getStockMaterialPackages(mergedFilter.value),
	enabled: computed(() => open),
})

function loadMore() {
	if (!stockMaterialPackages.value) return
	const pagination = stockMaterialPackages.value.pagination
	if (pagination.pageSize < pagination.totalCount) {
		localFilter.pageSize = (localFilter.pageSize ?? 10) + 10
	}
}

function selectMaterial(stockMaterial: StockMaterialPackagesDTO) {
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
	if (!stockMaterialPackages.value) return false
	const { pageSize, totalCount } = stockMaterialPackages.value.pagination
	return pageSize < totalCount
})
</script>

<style scoped></style>
