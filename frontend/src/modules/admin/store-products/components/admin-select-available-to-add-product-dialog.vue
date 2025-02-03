<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent
			:include-close-button="false"
			class="w-full max-w-2xl lg:max-w-4xl"
		>
			<DialogHeader>
				<DialogTitle>Выберите доступный для добавления товар</DialogTitle>
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
						v-if="!products || products.data.length === 0"
						class="text-muted-foreground"
					>
						Товары не найдены
					</p>

					<ul v-else>
						<li
							v-for="product in products.data"
							:key="product.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectProducts(product)"
						>
							<div class="flex items-center gap-2">
								<img
									:src="product.imageUrl"
									class="bg-gray-100 p-1 rounded-md w-16 h-16 object-contain"
								/>
								<span>{{ product.name }}</span>
							</div>
							<span class="text-gray-500 text-sm">
								{{ product.category.name }}
							</span>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="products && products.pagination.pageSize < products.pagination.totalCount"
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

import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import type { ProductDetailsDTO, ProductsFilterDTO } from '@/modules/kiosk/products/models/product.model'

const {open} = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', product: ProductDetailsDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<ProductsFilterDTO>({
  page: 1,
  pageSize: 10,
  search: ''
})


watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
})

const { data: products } = useQuery({
  queryKey: computed(() => [
  'admin-available-to-add-products',
  filter.value
]),
  queryFn: () => storeProductsService.getAvailableToAddProductsList(filter.value),
})


function loadMore() {
  if (!products.value) return
  const pagination = products.value.pagination

  if (pagination.pageSize < pagination.totalCount) {
    if(filter.value.pageSize) filter.value.pageSize += 10
  }
}

function selectProducts(product: ProductDetailsDTO) {
  emit('select', product)
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
