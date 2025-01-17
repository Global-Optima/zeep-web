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

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Добавить ингредиенты на склад</CardTitle>
						<CardDescription class="mt-2">
							Заполните таблицу ниже, чтобы добавить несколько ингредиентов на склад.
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
							<TableHead>Категория</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead>Порог малого запаса</TableHead>
							<TableHead class="text-center">Действия</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow v-if="selectedIngredients.length === 0">
							<TableCell
								colspan="4"
								class="py-5 text-center text-gray-500"
							>
								Нет добавленных ингредиентов
							</TableCell>
						</TableRow>
						<TableRow
							v-for="(ingredient, index) in selectedIngredients"
							:key="ingredient.ingredientId"
						>
							<TableCell>{{ ingredient.name }}</TableCell>
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
							<TableCell>
								<Input
									type="number"
									v-model.number="ingredient.lowStockThreshold"
									:min="0"
									:class="{ 'border-red-500': hasError(ingredient, 'lowStockThreshold') }"
									placeholder="Введите порог"
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
		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				>Сохранить</Button
			>
		</div>
	</div>

	<AdminIngredientsSelectDialog
		:open="openDialog"
		@close="openDialog = false"
		@select="addSelectedIngredient"
	/>
</template>

<script setup lang="ts">
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
  TableRow,
} from '@/core/components/ui/table'
import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import type { IngredientCategoryDTO, IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { AddMultipleStoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'
interface CreateStoreStockItem {
	name: string
	ingredientId: number
	quantity: number
	lowStockThreshold: number
  unit: UnitDTO
  category: IngredientCategoryDTO
}

// Props
const emit = defineEmits<{
	(e: 'onSubmit', payload: AddMultipleStoreWarehouseStockDTO): void
	(e: 'onCancel'): void
}>()

// State for selected ingredients
const selectedIngredients = ref<CreateStoreStockItem[]>([])

// Dialog open state
const openDialog = ref(false)

// Add selected ingredients to the table
function addSelectedIngredient(ingredient: IngredientsDTO) {
	const exists = selectedIngredients.value.some(
		(existingItem) => existingItem.ingredientId === ingredient.id
	)

	if (!exists) {
		selectedIngredients.value.push({
			name: ingredient.name,
			ingredientId: ingredient.id,
			quantity: 0,
			lowStockThreshold: 0,
      unit: ingredient.unit,
      category: ingredient.category
		})
		openDialog.value = false
	}
}

// Remove an ingredient from the table
function removeIngredient(index: number) {
	selectedIngredients.value.splice(index, 1)
}

// Check for validation errors
function hasError(ingredient: CreateStoreStockItem, field: 'quantity' | 'lowStockThreshold'): boolean {
	return ingredient[field] === undefined || ingredient[field] <= 0
}

// Computed: Can submit form
const canSubmit = computed(() => {
	return (
		selectedIngredients.value.length > 0 &&
		selectedIngredients.value.every((item) => !hasError(item, 'quantity') && !hasError(item, 'lowStockThreshold'))
	)
})

// Submit form
function onSubmit() {
	if (!canSubmit.value) {
		console.error('Validation errors detected.')
		return
	}

	const payload: AddMultipleStoreWarehouseStockDTO = {
		ingredientStocks: selectedIngredients.value.map((item) => ({
			ingredientId: item.ingredientId,
			quantity: item.quantity,
			lowStockThreshold: item.lowStockThreshold,
		})),
	}
	emit('onSubmit', payload)
}

// Cancel form
function onCancel() {
	emit('onCancel')
}
</script>
