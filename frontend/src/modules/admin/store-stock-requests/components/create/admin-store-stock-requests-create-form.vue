<template>
	<div>
		<Card>
			<CardHeader>
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Создать запрос на склад</CardTitle>
						<CardDescription class="mt-2">
							Заполните таблицу ниже, чтобы добавить материалы в запрос.
						</CardDescription>
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
							<TableCell>{{ item.name }}</TableCell>
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
			</CardContent>
			<CardFooter class="flex justify-end">
				<Button
					variant="outline"
					class="mr-2"
					@click="cancelForm"
				>
					Отмена
				</Button>
				<Button
					variant="default"
					@click="submitForm"
				>
					Сохранить
				</Button>
			</CardFooter>
		</Card>

		<AdminStockMaterialsSelectDialog
			:open="openDialog"
			@close="openDialog = false"
			@select="addMaterial"
		/>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { Trash } from 'lucide-vue-next'

import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'
import type { CreateStoreStockRequestItemDTO } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { ref } from 'vue'

interface StoreStockRequestItemForm extends CreateStoreStockRequestItemDTO {
  name: string
}

const emit = defineEmits<{
  (e: 'submit', payload: CreateStoreStockRequestItemDTO[]): void;
  (e: 'cancel'): void;
}>()

const stockRequestItemsForm = ref<StoreStockRequestItemForm[]>([])
const openDialog = ref(false)

function addMaterial(material: { id: number; name: string }) {
  if (!stockRequestItemsForm.value.some((item) => item.stockMaterialId === material.id)) {
    stockRequestItemsForm.value.push({
      stockMaterialId: material.id,
      name: material.name,
      quantity: 0,
    })
  }
}

function removeItem(index: number) {
  stockRequestItemsForm.value.splice(index, 1)
}

function submitForm() {
  emit('submit', stockRequestItemsForm.value)
}

function cancelForm() {
  emit('cancel')
}
</script>
