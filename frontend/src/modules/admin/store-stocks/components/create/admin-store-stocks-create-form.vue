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
							<TableHead class="text-center"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow v-if="selectedIngredients.length === 0">
							<TableCell
								colspan="5"
								class="py-5 text-gray-500 text-center"
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
  TableRow
} from '@/core/components/ui/table'
import { useToast } from '@/core/components/ui/toast'
import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import { ChevronLeft, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'

interface CreateStoreStockItem {
	name: string
	ingredientId: number
	quantity: number
	lowStockThreshold: number
	unit: { name: string }
	category: { name: string }
}

const emit = defineEmits<{
	(e: 'onSubmit', payload: { ingredientStocks: { ingredientId: number; quantity: number; lowStockThreshold: number }[] }): void
	(e: 'onCancel'): void
}>()

const selectedIngredients = ref<CreateStoreStockItem[]>([])
const openDialog = ref(false)
const { toast } = useToast()

// Add selected ingredients to the table
function addSelectedIngredient(ingredient: { id: number; name: string; unit: { name: string }; category: { name: string } }) {
	const exists = selectedIngredients.value.some(
		(existingItem) => existingItem.ingredientId === ingredient.id
	)

	if (exists) {
		toast({
			title: 'Ошибка',
			description: `Ингредиент "${ingredient.name}" уже добавлен.`,
			variant: 'destructive'
		})
		return
	}

	selectedIngredients.value.push({
		name: ingredient.name,
		ingredientId: ingredient.id,
		quantity: 0,
		lowStockThreshold: 0,
		unit: ingredient.unit,
		category: ingredient.category
	})

	toast({
		title: 'Успех',
		description: `Ингредиент "${ingredient.name}" добавлен.`,
		variant: 'default'
	})

	openDialog.value = false
}

// Remove ingredient from the table
function removeIngredient(index: number) {
	const removed = selectedIngredients.value.splice(index, 1)
	toast({
		title: 'Удалено',
		description: `Ингредиент "${removed[0].name}" удален.`,
		variant: 'default'
	})
}

// Validation checks
function hasError(item: CreateStoreStockItem, field: 'quantity' | 'lowStockThreshold'): boolean {
	return item[field] === undefined || item[field] <= 0
}

// Computed: Can submit form
const canSubmit = computed(() => {
	return (
		selectedIngredients.value.length > 0 &&
		selectedIngredients.value.every(
			(item) => !hasError(item, 'quantity') && !hasError(item, 'lowStockThreshold')
		)
	)
})

// Submit form
function onSubmit() {
	if (!canSubmit.value) {
		toast({
			title: 'Ошибка',
			description: 'Убедитесь, что все поля заполнены корректно.',
			variant: 'destructive'
		})
		return
	}

	const payload = {
		ingredientStocks: selectedIngredients.value.map((item) => ({
			ingredientId: item.ingredientId,
			quantity: item.quantity,
			lowStockThreshold: item.lowStockThreshold
		}))
	}

	toast({
		title: 'Успех',
		description: 'Ингредиенты успешно добавлены.',
		variant: 'default'
	})

	emit('onSubmit', payload)
}

// Cancel form
function onCancel() {
	toast({
		title: 'Отмена',
		description: 'Изменения отменены.',
		variant: 'default'
	})
	emit('onCancel')
}
</script>
