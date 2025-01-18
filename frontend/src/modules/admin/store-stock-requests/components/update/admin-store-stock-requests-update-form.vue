<template>
	<div>
		<Card>
			<CardHeader>
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Обновить запрос на склад</CardTitle>
						<CardDescription class="mt-2"> Отредактируйте материалы ниже. </CardDescription>
					</div>
					<Button
						variant="outline"
						@click="openDialog = true"
					>
						Добавить материал
					</Button>
				</div>
			</CardHeader>

			<CardContent>
				<form @submit.prevent="submitForm">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Материал</TableHead>
								<TableHead>Количество</TableHead>
								<TableHead class="text-center">Действия</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="(item, index) in stockRequestItemsForm"
								:key="index"
							>
								<TableCell>
									{{ item.name }}
								</TableCell>
								<TableCell>
									<Input
										type="number"
										v-model.number="item.quantity"
										class="p-1 w-full"
										placeholder="Введите количество"
									/>
								</TableCell>
								<TableCell class="flex justify-center text-center">
									<Trash
										class="text-red-500 hover:text-red-700 cursor-pointer"
										@click="removeItem(index)"
									/>
								</TableCell>
							</TableRow>
						</TableBody>
					</Table>

					<div class="flex justify-end mt-4">
						<Button
							variant="outline"
							class="mr-2"
							@click="cancelForm"
						>
							Отмена
						</Button>
						<Button
							variant="default"
							type="submit"
						>
							Сохранить
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>

		<AdminStockMaterialsSelectDialog
			:open="openDialog"
			@close="openDialog = false"
			@select="addMaterial"
		/>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'
import { Trash } from 'lucide-vue-next'

import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'
import type { StockRequestStockMaterialDTO, StockRequestStockMaterialResponse } from '@/modules/admin/store-stock-requests/models/stock-requests.model'

export interface StockRequestItemForm extends StockRequestStockMaterialDTO {
  name: string
}

const props = defineProps<{
  initialData: StockRequestStockMaterialResponse[]
}>()

const emit = defineEmits<{
  (e: 'submit', payload: StockRequestItemForm[]): void
  (e: 'cancel'): void
}>()

const stockRequestItemsForm = ref<StockRequestItemForm[]>([])
const openDialog = ref(false)

onMounted(() => {
  stockRequestItemsForm.value = props.initialData.map(item => ({ stockMaterialId: item.stockMaterialId,
    quantity:item.packageMeasures.quantity, name: item.name  }))
})

function addMaterial(material: { id: number; name: string }) {
  const itemExists = stockRequestItemsForm.value.some(
    item => item.stockMaterialId === material.id
  )
  if (!itemExists) {
    stockRequestItemsForm.value.push({
      stockMaterialId: material.id,
      name: material.name,
      quantity: 1
    })
  }
}

function removeItem(index: number) {
  stockRequestItemsForm.value.splice(index, 1)
}

function submitForm() {
  for (const item of stockRequestItemsForm.value) {
    if (!item.quantity || item.quantity <= 0) {
      alert(`Количество для "${item.name}" должно быть больше нуля!`)
      return
    }
  }
  emit('submit', stockRequestItemsForm.value)
}

function cancelForm() {
  emit('cancel')
}
</script>
