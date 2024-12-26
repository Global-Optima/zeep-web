<template>
	<div>
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
				<!-- Table for Batch Ingredient Input -->
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Ингредиент</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead>Порог малого запаса</TableHead>
							<TableHead class="text-center">Действия</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="(ingredient, index) in selectedIngredients"
							:key="index"
						>
							<TableCell>{{ ingredient.name }}</TableCell>
							<TableCell>
								<Input
									type="number"
									v-model.number="ingredient.quantity"
									class="p-1 w-full"
									placeholder="Введите количество"
								/>
							</TableCell>
							<TableCell>
								<Input
									type="number"
									v-model.number="ingredient.lowStockThreshold"
									class="p-1 w-full"
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

		<!-- Ingredient Selection Dialog -->
		<Dialog v-model:open="openDialog">
			<DialogContent :includeCloseButton="false">
				<DialogHeader>
					<DialogTitle>Выберите ингредиент</DialogTitle>
				</DialogHeader>
				<div>
					<!-- Search Input -->
					<Input
						v-model="searchTerm"
						placeholder="Поиск ингредиента"
						type="search"
						class="mt-2 mb-4 w-full"
					/>

					<!-- Ingredient List -->
					<div class="max-h-64 overflow-y-auto">
						<ul>
							<li
								v-for="ingredient in filteredIngredients"
								:key="ingredient.id"
								class="flex justify-between items-center hover:bg-gray-100 p-2 border-b cursor-pointer"
								@click="addIngredient(ingredient)"
							>
								<span>{{ ingredient.name }}</span>
							</li>
						</ul>
					</div>
				</div>
				<DialogFooter>
					<Button
						variant="outline"
						@click="openDialog = false"
						>Закрыть</Button
					>
				</DialogFooter>
			</DialogContent>
		</Dialog>
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
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
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
import { computed, ref } from 'vue'

interface CreateStoreStockItem {
	ingredientId: number
	quantity: number
	lowStockThreshold: number
}

export interface CreateMultipleStoreStock {
	ingredientStocks: CreateStoreStockItem[]
}

interface Ingredient {
	id: number
	name: string
}

interface SelectedIngredient extends CreateStoreStockItem {
	name: string
}

// Props for parent component
const emit = defineEmits<{
	(e: 'onSubmit', payload: CreateMultipleStoreStock): void
	(e: 'onCancel'): void
}>()

// State for selected ingredients
const selectedIngredients = ref<SelectedIngredient[]>([])

// Dialog open state
const openDialog = ref(false)

// Ingredient list (mocked for now, should come from an API)
const allIngredients = ref<Ingredient[]>([
	{ id: 1, name: 'Сахар' },
	{ id: 2, name: 'Мука' },
	{ id: 3, name: 'Соль' },
])

// Search term for filtering ingredients
const searchTerm = ref('')

// Filtered ingredient list based on search term
const filteredIngredients = computed(() =>
	allIngredients.value.filter((ingredient) =>
		ingredient.name.toLowerCase().includes(searchTerm.value.toLowerCase())
	)
)

// Add ingredient to the selected list
function addIngredient(ingredient: Ingredient) {
	if (!selectedIngredients.value.some((i) => i.ingredientId === ingredient.id)) {
		selectedIngredients.value.push({
			ingredientId: ingredient.id,
			name: ingredient.name,
			quantity: 0,
			lowStockThreshold: 0,
		})
	}
	openDialog.value = false
}

// Remove ingredient from the selected list
function removeIngredient(index: number) {
	selectedIngredients.value.splice(index, 1)
}

// Submit form
function submitForm() {
	const payload: CreateMultipleStoreStock = {
		ingredientStocks: selectedIngredients.value.map((item) => ({
			ingredientId: item.ingredientId,
			quantity: item.quantity,
			lowStockThreshold: item.lowStockThreshold,
		})),
	}
	emit('onSubmit', payload)
}

// Cancel form
function cancelForm() {
	emit('onCancel')
}
</script>
