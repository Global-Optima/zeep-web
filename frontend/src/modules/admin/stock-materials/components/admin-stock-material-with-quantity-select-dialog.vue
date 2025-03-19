<template>
	<Dialog
		:open="open"
		@update:open="onClose"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Добавить материал в заявку</DialogTitle>
			</DialogHeader>

			<div>
				<!-- Selected Stock Material -->
				<div class="mb-4">
					<Label>Материал</Label>
					<p
						class="mt-1 text-primary underline cursor-pointer"
						@click="openMaterialDialog"
					>
						{{ selectedStockMaterialLabel }}
					</p>
				</div>

				<!-- Quantity Input -->
				<div class="mb-4">
					<Label for="quantity">Количество упаковок</Label>
					<Input
						id="quantity"
						type="number"
						min="1"
						v-model.number="quantity"
						placeholder="Введите количество"
						class="mt-2"
					/>
					<p
						v-if="quantity <= 0"
						class="mt-1 text-red-500 text-sm"
					>
						Количество должно быть больше 0
					</p>
				</div>
			</div>

			<!-- Dialog Footer -->
			<DialogFooter>
				<Button
					variant="outline"
					@click="onClose"
					>Закрыть</Button
				>
				<Button
					variant="default"
					:disabled="!selectedStockMaterial || quantity <= 0"
					@click="handleSubmit"
				>
					Добавить
				</Button>
			</DialogFooter>
		</DialogContent>

		<!-- Select Material Dialog -->
		<AdminSelectAvailableToAddStockMaterialsDialog
			:open="materialDialogOpen"
			:initial-filter="initialFilter"
			@close="materialDialogOpen = false"
			@select="handleMaterialSelect"
		/>
	</Dialog>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import type { StockMaterialsDTO, StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StockRequestStockMaterialDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
import AdminSelectAvailableToAddStockMaterialsDialog from '@/modules/admin/warehouse-stocks/components/admin-select-available-to-add-stock-materials-dialog.vue'
import { computed, ref } from 'vue'

const { open, initialFilter } = defineProps<{ open: boolean, initialFilter?: StockMaterialsFilter}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', dto: StockRequestStockMaterialDTO): void
}>()

const selectedStockMaterial = ref<StockMaterialsDTO | null>(null)
const quantity = ref<number>(1)
const materialDialogOpen = ref(false)

const selectedStockMaterialLabel = computed(() => {
  const material = selectedStockMaterial.value

  return material ?
  `${material.name}, ${material.size} ${material.unit.name.toLowerCase()}`
  : 'Выбрать материал'
})

function openMaterialDialog() {
  materialDialogOpen.value = true
}

function handleMaterialSelect(stockMaterial: StockMaterialsDTO) {
  selectedStockMaterial.value = stockMaterial
  materialDialogOpen.value = false
}

/** Close the dialog */
function onClose() {
  selectedStockMaterial.value = null
  quantity.value = 1
  emit('close')
}

/** Submit the form */
function handleSubmit() {
  if (selectedStockMaterial.value && quantity.value > 0) {
    emit('submit', {
      stockMaterialId: selectedStockMaterial.value.id,
      quantity: quantity.value
    })
    onClose()
  }
}
</script>

<style scoped>
/* Scoped styles for better accessibility and UI */
</style>
