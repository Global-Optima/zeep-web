<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Выберите заготовку</DialogTitle>
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
						v-if="!provisions || provisions.data.length === 0"
						class="text-muted-foreground"
					>
						Заготовки не найдены
					</p>

					<ul v-else>
						<li
							v-for="provision in provisions.data"
							:key="provision.id"
							class="flex justify-between items-start hover:bg-gray-100 px-2 py-3 border-b rounded-lg cursor-pointer"
							@click="selectMaterial(provision)"
						>
							<span>{{ provision.name }}</span>

							<span class="text-gray-600"
								>{{ provision.absoluteVolume }} {{ provision.unit.name.toLowerCase() }}</span
							>
						</li>
					</ul>
				</div>

				<!-- Load More Button -->
				<Button
					v-if="provisions && provisions.pagination.pageSize < provisions.pagination.totalCount"
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

import type { ProvisionDTO, ProvisionFilter } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"

const {open} = defineProps<{
  open: boolean;
}>()

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', provision: ProvisionDTO): void;
}>()

const searchTerm = ref('')
const debouncedSearchTerm = useDebounce(
  computed(() => searchTerm.value),
  500
)

const filter = ref<ProvisionFilter>({})


watch(debouncedSearchTerm, (newValue) => {
  filter.value.page = 1
  filter.value.search = newValue.trim()
})

const { data: provisions } = useQuery({
  queryKey: computed(() => [
  'admin-provisions',
  filter.value
]),
  queryFn: () => provisionsService.getProvisions(filter.value),
})


function loadMore() {
  if (!provisions.value) return
  const pagination = provisions.value.pagination

  if (pagination.pageSize < pagination.totalCount) {
    if(filter.value.pageSize) filter.value.pageSize += 10
  }
}

function selectMaterial(provision: ProvisionDTO) {
  emit('select', provision)
  onClose()
}

function onClose() {
  filter.value = {}
  emit('close')
}
</script>

<style scoped></style>
