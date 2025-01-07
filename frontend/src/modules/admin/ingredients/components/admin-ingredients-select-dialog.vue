<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите ингредиент</DialogTitle>
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
						v-if="!ingredients || ingredients.data.length === 0"
						class="text-muted-foreground"
					>
						Ингредиенты не найдены
					</p>

					<ul v-else>
						<li
							v-for="ingredient in ingredients.data"
							:key="ingredient.id"
							class="flex justify-between items-center hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(ingredient)"
						>
							<span>{{ ingredient.name }}</span>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="ingredients && ingredients.pagination.page < ingredients.pagination.totalPages"
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

import type { IngredientFilter, IngredientResponseDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'

const {open} = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', additive: IngredientResponseDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<IngredientFilter>({
  page: 1,
  pageSize: 10,
  name: ''
})


watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.name = newValue.trim()
  refetch()
})

const { data: ingredients, refetch } = useQuery({
  queryKey: computed(() => [
  'admin-ingredients',
  filter.value
]),
  queryFn: () => ingredientsService.getIngredients(filter.value),
})


function loadMore() {
  if (!ingredients.value) return
  const pagination = ingredients.value.pagination

  if (pagination.pageSize < pagination.totalCount) {
    if(filter.value.pageSize) filter.value.pageSize += 10
  }
}

function selectMaterial(ing: IngredientResponseDTO) {
  emit('select', ing)
  onClose()
}

function onClose() {
  filter.value = {
    page: 1,
    pageSize: 10,
    name: ''
  }

  emit('close')
}
</script>

<style scoped></style>
