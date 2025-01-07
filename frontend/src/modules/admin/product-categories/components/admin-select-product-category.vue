<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import type { ProductCategoriesFilterDTO, ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

const { open } = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', additiveCategory: ProductCategoryDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<ProductCategoriesFilterDTO>({
  page: 1,
  pageSize: 10,
  search: ''
})

watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
  refetch()
})

const { data: additiveCategories, refetch } = useQuery({
  queryKey: computed(() => [
    'admin-products-categories',
    filter.value
  ]),
  queryFn: () => productsService.getAllProductCategories(filter.value),
})

function loadMore() {
  if (!additiveCategories.value) return
  const pagination = additiveCategories.value.pagination
  if (filter.value.pageSize && pagination.pageSize < pagination.totalCount) {
    filter.value.pageSize += 10
    refetch()
  }
}

function selectAdditiveCategory(productCategory: ProductCategoryDTO) {
  emit('select', productCategory)
  onClose()
}

function handleDialogState(newState: boolean) {
  if (!newState) onClose()
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

<template>
	<Dialog
		:open="open"
		@update:open="handleDialogState"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите категорию товара</DialogTitle>
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
						v-if="!additiveCategories || additiveCategories.data.length === 0"
						class="text-muted-foreground"
					>
						Категории товаров не найдены
					</p>
					<ul v-else>
						<li
							v-for="additiveCategory in additiveCategories.data"
							:key="additiveCategory.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectAdditiveCategory(additiveCategory)"
						>
							<span>{{ additiveCategory.name }}</span>
						</li>
					</ul>
				</div>

				<Button
					v-if="additiveCategories && additiveCategories.pagination.pageSize < additiveCategories.pagination.totalCount"
					variant="ghost"
					type="button"
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
					type="button"
				>
					Закрыть
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
