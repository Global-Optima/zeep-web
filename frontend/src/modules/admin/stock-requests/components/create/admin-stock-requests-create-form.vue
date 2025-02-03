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

import { useToast } from '@/core/components/ui/toast'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StockRequestStockMaterialDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import AdminSelectAvailableToAddStockMaterialsDialog from '@/modules/admin/warehouse-stocks/components/admin-select-available-to-add-stock-materials-dialog.vue'
import { ref } from 'vue'

interface CreateStockRequestItemForm {
	stockMaterialId: number
	name: string
	quantity: number
  size: number
  unit: UnitDTO
}

const emit = defineEmits<{
  (e: 'submit', payload: StockRequestStockMaterialDTO[]): void;
  (e: 'cancel'): void;
}>()

const stockRequestItemsForm = ref<CreateStockRequestItemForm[]>([])
const openDialog = ref(false)
const { toast } = useToast()

// Add Material
function addMaterial(material: StockMaterialsDTO) {
  if (stockRequestItemsForm.value.some((item) => item.stockMaterialId === material.id)) {
    toast({
      title: 'Ошибка',
      description: 'Этот материал уже добавлен.',
      variant: 'destructive',
    })
    return
  }

  stockRequestItemsForm.value.push({
    stockMaterialId: material.id,
    name: material.name,
    quantity: 0,
    unit: material.unit,
    size: material.size,
  })
  toast({
    title: 'Успех',
    description: `Материал "${material.name}" добавлен.`,
    variant: 'default',
  })
}

// Remove Material
function removeItem(index: number) {
  const removedItem = stockRequestItemsForm.value.splice(index, 1)[0]
  toast({
    title: 'Удалено',
    description: `Материал "${removedItem.name}" удален.`,
    variant: 'default',
  })
}

// Submit Form with Validations
function submitForm() {
  if (stockRequestItemsForm.value.length === 0) {
    toast({
      title: 'Ошибка',
      description: 'Добавьте хотя бы один материал перед отправкой.',
      variant: 'destructive',
    })
    return
  }

  const invalidItems = stockRequestItemsForm.value.filter((item) => item.quantity <= 0)
  if (invalidItems.length > 0) {
    toast({
      title: 'Ошибка',
      description: 'Убедитесь, что все количества больше нуля.',
      variant: 'destructive',
    })
    return
  }

  toast({
    title: 'Успех',
    description: 'Запрос успешно отправлен.',
    variant: 'default',
  })

  emit('submit', stockRequestItemsForm.value)
}

// Cancel Form
function cancelForm() {
  toast({
    title: 'Отмена',
    description: 'Действие отменено.',
    variant: 'default',
  })
  emit('cancel')
}
</script>

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
							<TableHead>Упаковка</TableHead>
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
							<TableCell>{{ item.size }} {{ item.unit.name }}</TableCell>
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

		<AdminSelectAvailableToAddStockMaterialsDialog
			:open="openDialog"
			@close="openDialog = false"
			@select="addMaterial"
		/>
	</div>
</template>
