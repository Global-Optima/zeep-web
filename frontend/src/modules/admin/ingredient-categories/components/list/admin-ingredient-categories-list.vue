<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="category in ingredientCategories"
				:key="category.id"
				class="h-12 cursor-pointer"
				@click="goToDetails(category.id)"
			>
				<TableCell class="p-4">
					<span class="font-medium">{{ category.name }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span class="max-w-12 font-medium truncate">{{ category.description }}</span>
				</TableCell>
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
import type { IngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { useRouter } from 'vue-router'

const {ingredientCategories} = defineProps<{ingredientCategories: IngredientCategoryDTO[]}>()

const router = useRouter();

const goToDetails = (ingredientCategoryId: number) => {
  router.push(`/admin/ingredient-categories/${ingredientCategoryId}`);
};
</script>
