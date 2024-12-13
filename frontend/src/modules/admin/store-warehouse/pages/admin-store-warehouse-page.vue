<template>
	<div>
		<!-- If no ingredients, display message -->
		<p
			v-if="ingredients.data.length === 0"
			class="text-muted-foreground"
		>
			Ингредиенты не найдены
		</p>
		<!-- If there are ingredients, display the table -->
		<Table v-else>
			<TableHeader>
				<TableRow>
					<TableHead>Название</TableHead>
					<TableHead>Категория</TableHead>
					<TableHead>Количество</TableHead>
					<TableHead class="hidden md:table-cell">Единица измерения</TableHead>
					<TableHead class="hidden md:table-cell">Статус</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<TableRow
					v-for="ingredient in displayedIngredients"
					:key="ingredient.id"
					class="hover:bg-gray-50 h-12 cursor-pointer"
				>
					<TableCell class="font-medium">
						{{ ingredient.name }}
					</TableCell>
					<TableCell class="font-medium">
						{{ ingredient.category }}
					</TableCell>
					<TableCell class="font-medium">
						{{ ingredient.quantity }}
					</TableCell>
					<TableCell class="hidden font-medium md:table-cell">
						{{ ingredient.unit }}
					</TableCell>
					<TableCell class="hidden font-medium md:table-cell">
						<p
							:class="[ 
								'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
								INGREDIENT_STATUS_COLOR[ingredient.status]
							]"
						>
							{{ INGREDIENT_STATUS_FORMATTED[ingredient.status] }}
						</p>
					</TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
</template>

<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { computed, ref } from 'vue'

// Constants
const DEFAULT_LIMIT = 10

// Mock data for warehouse ingredients
const ingredients = ref({
  data: [
    {
      id: 1,
      name: 'Кофейные зерна',
      category: 'Кофе',
      quantity: 50,
      unit: 'кг',
      status: 'in_stock',
    },
    {
      id: 2,
      name: 'Молоко',
      category: 'Молочные продукты',
      quantity: 100,
      unit: 'литр',
      status: 'in_stock',
    },
    {
      id: 3,
      name: 'Сахар',
      category: 'Сыпучие',
      quantity: 30,
      unit: 'кг',
      status: 'low_stock',
    },
    {
      id: 4,
      name: 'Сироп ванильный',
      category: 'Сиропы',
      quantity: 20,
      unit: 'литр',
      status: 'in_stock',
    },
    {
      id: 5,
      name: 'Чай черный',
      category: 'Чай',
      quantity: 40,
      unit: 'кг',
      status: 'in_stock',
    },
    {
      id: 6,
      name: 'Чай зеленый',
      category: 'Чай',
      quantity: 10,
      unit: 'кг',
      status: 'low_stock',
    },
    {
      id: 7,
      name: 'Сироп клубничный',
      category: 'Сиропы',
      quantity: 5,
      unit: 'литр',
      status: 'out_of_stock',
    },
    {
      id: 8,
      name: 'Шоколадный порошок',
      category: 'Сыпучие',
      quantity: 25,
      unit: 'кг',
      status: 'in_stock',
    },
    {
      id: 9,
      name: 'Корица',
      category: 'Специи',
      quantity: 3,
      unit: 'кг',
      status: 'low_stock',
    },
    {
      id: 10,
      name: 'Взбитые сливки',
      category: 'Молочные продукты',
      quantity: 0,
      unit: 'литр',
      status: 'out_of_stock',
    },
  ],
  meta: {
    totalItems: 10, // Total number of ingredients (for pagination)
  },
})

const limit = ref(DEFAULT_LIMIT)

// Computed property for displayed ingredients based on limit
const displayedIngredients = computed(() =>
  ingredients.value.data.slice(0, limit.value)
)

// Ingredient status colors and formatted text
const INGREDIENT_STATUS_COLOR: Record<string, string> = {
  in_stock: 'bg-green-100 text-green-800',
  low_stock: 'bg-yellow-100 text-yellow-800',
  out_of_stock: 'bg-red-100 text-red-800',
}

const INGREDIENT_STATUS_FORMATTED: Record<string, string> = {
  in_stock: 'В наличии',
  low_stock: 'Заканчивается',
  out_of_stock: 'Нет в наличии',
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
