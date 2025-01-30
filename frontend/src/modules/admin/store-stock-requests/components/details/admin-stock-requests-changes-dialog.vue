<template>
	<Dialog
		:open="open"
		@update:open="$emit('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="w-full max-w-3xl"
		>
			<DialogHeader>
				<DialogTitle>Принять с изменениями</DialogTitle>
				<DialogDescription>
					Укажите итоговое количество для каждого материала и добавьте комментарий, если нужно.
				</DialogDescription>
			</DialogHeader>

			<!-- COMMENT INPUT -->
			<div class="mb-4">
				<Label for="comment">Комментарий</Label>
				<Textarea
					id="comment"
					v-model="comment"
					placeholder="Оставьте комментарий"
					rows="2"
					class="mt-1"
				/>
			</div>

			<!-- MATERIALS TABLE -->
			<div class="overflow-auto">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Материал</TableHead>
							<TableHead>Ожидалось</TableHead>
							<TableHead>Фактически</TableHead>
							<TableHead class="text-center"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="(item, idx) in editableItems"
							:key="item.stockMaterialId"
						>
							<TableCell>{{ item.name }}</TableCell>
							<TableCell>{{ item.expectedQuantity }}</TableCell>
							<TableCell>
								<Input
									type="number"
									class="w-24"
									:class="{'border-red-500': item.quantity < 0}"
									v-model.number="item.quantity"
									@change="validateQuantity(idx)"
									min="0"
								/>
							</TableCell>
							<TableCell class="text-center">
								<Button
									variant="ghost"
									size="icon"
									class="text-red-500 hover:text-red-700"
									@click="removeItem(idx)"
								>
									<Trash class="w-5 h-5" />
								</Button>
							</TableCell>
						</TableRow>
						<!-- If no items -->
						<TableRow v-if="editableItems.length === 0">
							<TableCell
								colspan="4"
								class="py-4 text-center text-muted-foreground"
							>
								Нет материалов
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</div>

			<!-- ADD MATERIAL BUTTON -->
			<div class="flex justify-center mt-4">
				<Button
					variant="outline"
					@click="openSelectDialog = true"
				>
					Добавить материал
				</Button>
			</div>

			<DialogFooter>
				<Button
					variant="outline"
					@click="$emit('close')"
				>
					Отмена
				</Button>
				<Button
					variant="default"
					@click="submitForm"
				>
					Принять
				</Button>
			</DialogFooter>
		</DialogContent>

		<!-- MATERIAL SELECT DIALOG (for choosing new materials to add) -->
		<AdminStockMaterialsSelectDialog
			:open="openSelectDialog"
			@close="openSelectDialog = false"
			@select="addMaterial"
		/>
	</Dialog>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Dialog, DialogContent, DialogDescription,
  DialogFooter, DialogHeader, DialogTitle
} from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import {
  Table, TableBody, TableCell,
  TableHead, TableHeader, TableRow
} from '@/core/components/ui/table'
import { Textarea } from '@/core/components/ui/textarea'
import { Trash } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'

// The final DTO shape for submission
import type {
  AcceptWithChangeRequestStatusDTO
} from '@/modules/admin/store-stock-requests/models/stock-requests.model'

// The raw data we get for each item from the backend
import type { StockRequestMaterial } from '@/modules/admin/store-stock-requests/models/stock-requests.model'


/**
 * We'll define a simpler local type for the row in the table.
 * This is what the component will actually display and edit.
 * We transform from "StockRequestMaterial" -> "EditableMaterialRow".
 */
interface EditableMaterialRow {
  stockMaterialId: number
  name: string
  expectedQuantity: number
  quantity: number
}

const props = defineProps<{
  open: boolean
  initialItems: StockRequestMaterial[]
}>()

const emit = defineEmits<{
  (e: 'close'): void,
  (e: 'submit', payload: AcceptWithChangeRequestStatusDTO): void
}>()

/**
 * We'll store a local reactive array of "EditableMaterialRow".
 */
const editableItems = ref<EditableMaterialRow[]>([])

/**
 * Convert the incoming StockRequestMaterial[] into simpler EditableMaterialRow[].
 */


// The comment the user can provide
const comment = ref('')
const openSelectDialog = ref(false)

/** Initialize / watch for changes in props.initialItems. */
watch(
  () => props.initialItems,
  (newVal) => {
    editableItems.value = newVal.map(toEditableRow)
    comment.value = '' // Reset comment if needed
  },
  { immediate: true, deep: true }
)

function toEditableRow(item: StockRequestMaterial): EditableMaterialRow {
  return {
    stockMaterialId: item.stockMaterial.id,
    name: item.stockMaterial.name,
    expectedQuantity: item.quantity,
    quantity: item.quantity
  }
}

/**
 * Called when the user picks a new material from AdminStockMaterialsSelectDialog.
 * We'll add it to editableItems with default expected=0, quantity=0.
 */
function addMaterial(material: { id: number; name: string }) {
  openSelectDialog.value = false
  if (editableItems.value.some(item => item.stockMaterialId === material.id)) return

  editableItems.value.push({
    stockMaterialId: material.id,
    name: material.name,
    expectedQuantity: 0, // if you have logic for "expected," set it
    quantity: 0
  })
}

/** Remove a row. */
function removeItem(index: number) {
  editableItems.value.splice(index, 1)
}

/** Validate: ensure quantity >= 0. */
function validateQuantity(index: number) {
  const item = editableItems.value[index]
  if (item.quantity < 0) {
    item.quantity = 0
  }
}

/**
 * On submit: we build the final DTO:
 * {
 *   comment: string,
 *   items: Array<{ stockMaterialId, quantity }>
 * }
 */
function submitForm() {
  // Ensure no negative values
  for (const row of editableItems.value) {
    if (row.quantity < 0) {
      row.quantity = 0
    }
  }

  const dto: AcceptWithChangeRequestStatusDTO = {
    comment: comment.value.trim(),
    items: editableItems.value.map((r) => ({
      stockMaterialId: r.stockMaterialId,
      quantity: r.quantity
    }))
  }

  emit('submit', dto)
}
</script>

<style scoped>
/* Optional styling overrides */
</style>
