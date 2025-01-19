<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import type { SuppliersDTO, SuppliersFilter } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

const { open } = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', category: SuppliersDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<SuppliersFilter>({})

watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.pageSize = 10
  filter.value.search = newValue.trim()
})

const { data: categories} = useQuery({
  queryKey: computed(() => [
    'admin-suppliers',
    filter.value
  ]),
  queryFn: () => suppliersService.getSuppliers(filter.value),
})

function loadMore() {
  if (!categories.value) return
  const pagination = categories.value.pagination
  if (filter.value.pageSize && pagination.pageSize < pagination.totalCount) {
    filter.value.pageSize += 10
  }
}

function selectCategory(productCategory: SuppliersDTO) {
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
  searchTerm.value = ""
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
				<DialogTitle>Выберите поставщика товаров</DialogTitle>
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
						v-if="!categories || categories.data.length === 0"
						class="text-muted-foreground"
					>
						Поставщики не найдены
					</p>
					<ul v-else>
						<li
							v-for="category in categories.data"
							:key="category.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectCategory(category)"
						>
							<span>{{ category.name }}</span>
						</li>
					</ul>
				</div>

				<Button
					v-if="categories && categories.pagination.pageSize < categories.pagination.totalCount"
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
