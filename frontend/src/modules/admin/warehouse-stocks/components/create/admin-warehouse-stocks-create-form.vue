<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Добавить на склад
			</h1>

			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Добавить продукты на склад</CardTitle>
						<CardDescription class="mt-2">
							Заполните таблицу ниже, чтобы добавить несколько продуктов на склад.
						</CardDescription>
					</div>
					<Button
						variant="outline"
						@click="openDialog = true"
					>
						Добавить
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Упаковка</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead class="text-center"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow v-if="selectedStockMaterials.length === 0">
							<TableCell
								colspan="5"
								class="py-5 text-gray-500 text-center"
							>
								Сырье не добавлено
							</TableCell>
						</TableRow>
						<TableRow
							v-for="(ingredient, index) in selectedStockMaterials"
							:key="ingredient.stockMaterialId"
						>
							<TableCell>{{ ingredient.name }}</TableCell>
							<TableCell>{{ ingredient.size }} {{ ingredient.unit.name }}</TableCell>
							<TableCell>{{ ingredient.category.name }}</TableCell>
							<TableCell>
								<Input
									type="number"
									v-model.number="ingredient.quantity"
									:min="0"
									:class="{ 'border-red-500': hasError(ingredient, 'quantity') }"
									placeholder="Введите количество"
								/>
							</TableCell>
							<TableCell class="flex justify-center text-center">
								<Trash
									class="text-red-500 hover:text-red-700 cursor-pointer"
									@click="removeIngredient(index)"
								/>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="onCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>

		<!-- Dialog -->
		<AdminStockMaterialsSelectDialog
			:open="openDialog"
			@close="openDialog = false"
			@select="addStockMaterial"
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
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'
// Interfaces
interface CreateWarehouseStockItem {
	name: string
	stockMaterialId: number
	quantity: number
	unit: UnitDTO
	category: StockMaterialCategoryDTO
	ingredient: IngredientsDTO
  size: number
}

const emit = defineEmits<{
	(e: 'onSubmit', payload: { stockMaterialId: number; quantity: number }[]): void
	(e: 'onCancel'): void
}>()

const selectedStockMaterials = ref<CreateWarehouseStockItem[]>([])
const openDialog = ref(false)
const { toast } = useToast()

// Add Material
function addStockMaterial(stockMaterial: StockMaterialsDTO) {
	const exists = selectedStockMaterials.value.some(
		(item) => item.stockMaterialId === stockMaterial.id
	)

	if (exists) {
		toast({
			title: 'Ошибка',
			description: `Материал "${stockMaterial.name}" уже добавлен.`,
			variant: 'destructive',
		})
		return
	}

	selectedStockMaterials.value.push({
    quantity: 0,
    name: stockMaterial.name,
    stockMaterialId: stockMaterial.id,
    unit: stockMaterial.unit,
    category: stockMaterial.category,
    ingredient: stockMaterial.ingredient,
    size: stockMaterial.size
  })

	toast({
		title: 'Успех',
		description: `Материал "${stockMaterial.name}" добавлен.`,
		variant: 'default',
	})
	openDialog.value = false
}

// Remove Material
function removeIngredient(index: number) {
	const removed = selectedStockMaterials.value.splice(index, 1)
	toast({
		title: 'Удалено',
		description: `Материал "${removed[0].name}" удален.`,
		variant: 'default',
	})
}

// Validation Checks
function hasError(item: CreateWarehouseStockItem, field: 'quantity'): boolean {
	return item[field] === undefined || item[field] < 0
}

// Enable Submit
const canSubmit = computed(() => {
	return (
		selectedStockMaterials.value.length > 0 &&
		selectedStockMaterials.value.every((item) => !hasError(item, 'quantity'))
	)
})

// Submit Form
function onSubmit() {
	if (!canSubmit.value) {
		toast({
			title: 'Ошибка',
			description: 'Убедитесь, что все поля заполнены корректно.',
			variant: 'destructive',
		})
		return
	}

	const payload = selectedStockMaterials.value.map((item) => ({
		stockMaterialId: item.stockMaterialId,
		quantity: item.quantity,
	}))

	toast({
		title: 'Успех',
		description: 'Материалы успешно добавлены.',
		variant: 'default',
	})

	emit('onSubmit', payload)
}

// Cancel
function onCancel() {
	emit('onCancel')
}
</script>
