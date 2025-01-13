<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Калорий (ккал)</TableHead>
				<TableHead>Белков (гр)</TableHead>
				<TableHead>Жиров (гр)</TableHead>
				<TableHead>Углеводов (гр)</TableHead>
				<TableHead>Срок годности (дни)</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="ingredient in ingredients"
				:key="ingredient.id"
				@click="onIngredientClick(ingredient.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ ingredient.name }}</TableCell>
				<TableCell>{{ ingredient.calories }}</TableCell>
				<TableCell>{{ ingredient.proteins }}</TableCell>
				<TableCell>{{ ingredient.fat }}</TableCell>
				<TableCell>{{ ingredient.carbs }}</TableCell>
				<TableCell
					>{{ ingredient.expirationInDays ?? "Не указано" }}</TableCell
				>
			</TableRow>
		</TableBody>
	</Table>
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
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { format } from 'date-fns'
import { useRouter } from 'vue-router'

const {ingredients} = defineProps<{ingredients: IngredientsDTO[]}>()

const router = useRouter();

const onIngredientClick = (ingredientId: number) => {
  router.push(`/admin/ingredients/${ingredientId}`);
};
</script>
