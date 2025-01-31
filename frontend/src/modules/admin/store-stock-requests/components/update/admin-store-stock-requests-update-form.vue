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

		<AdminSelectAvailableToAddStockMaterialsDialog
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
import { useToast } from '@/core/components/ui/toast'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StockRequestMaterial } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import AdminSelectAvailableToAddStockMaterialsDialog from '@/modules/admin/warehouse-stocks/components/admin-select-available-to-add-stock-materials-dialog.vue'
import { Trash } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
export interface StockRequestItemForm {
	stockMaterialId: number
	name: string
	quantity: number
  size: number
  unit: UnitDTO
}

const props = defineProps<{
	initialData: StockRequestMaterial[]
}>()

const emit = defineEmits<{
	(e: 'submit', payload: StockRequestItemForm[]): void
	(e: 'cancel'): void
}>()

const stockRequestItemsForm = ref<StockRequestItemForm[]>([])
const openDialog = ref(false)
const { toast } = useToast()

onMounted(() => {
	stockRequestItemsForm.value = props.initialData.map((item) => ({
		stockMaterialId: item.stockMaterial.id,
		name: item.stockMaterial.name,
		quantity: item.quantity,
    unit: item.stockMaterial.unit,
    size: item.stockMaterial.size,
	}))
})

function addMaterial(material: StockMaterialsDTO) {
	if (stockRequestItemsForm.value.some((item) => item.stockMaterialId === material.id)) {
		toast({
			title: 'Ошибка',
			description: `Материал "${material.name}" уже добавлен.`,
			variant: 'destructive',
		})
		return
	}

	stockRequestItemsForm.value.push({
		stockMaterialId: material.id,
		name: material.name,
		quantity: 1,
    unit: material.unit,
    size: material.size,
	})

	toast({
		title: 'Успех',
		description: `Материал "${material.name}" добавлен.`,
		variant: 'default',
	})
}

function removeItem(index: number) {
	const removedItem = stockRequestItemsForm.value.splice(index, 1)[0]
	toast({
		title: 'Удалено',
		description: `Материал "${removedItem.name}" удален.`,
		variant: 'default',
	})
}

function submitForm() {
	if (stockRequestItemsForm.value.length === 0) {
		toast({
			title: 'Ошибка',
			description: 'Добавьте хотя бы один материал перед сохранением.',
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
		description: 'Запрос успешно обновлен.',
		variant: 'default',
	})

	emit('submit', stockRequestItemsForm.value)
}

function cancelForm() {
	toast({
		title: 'Отмена',
		description: 'Изменения отменены.',
		variant: 'default',
	})
	emit('cancel')
}
</script>
