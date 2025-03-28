<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите размер</DialogTitle>
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
						v-if="!additiveCategories || additiveCategories.data.length === 0"
						class="text-muted-foreground"
					>
						Размер не найден
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
				<!-- Load More Button -->
				<Button
					v-if="additiveCategories && additiveCategories.pagination.pageSize < additiveCategories.pagination.totalCount"
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
import type { UnitDTO, UnitsFilterDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
const {open} = defineProps<{
  open: boolean;
}>()
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', unit: UnitDTO): void;
}>()
const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)
const filter = ref<UnitsFilterDTO>({})

watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
})

const { data: additiveCategories } = useQuery({
  queryKey: computed(() => [
  'admin-units',
  filter.value
]),
  queryFn: () => unitsService.getAllUnits(filter.value),
})

function loadMore() {
  if (!additiveCategories.value) return
  const pagination = additiveCategories.value.pagination
  if (pagination.pageSize < pagination.totalCount) {
    if(filter.value.pageSize) filter.value.pageSize += 10
  }
}

function selectAdditiveCategory(unit: UnitDTO) {
  emit('select', unit)
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
