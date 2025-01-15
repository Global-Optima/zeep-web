<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead class="p-4"></TableHead>
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

				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, category.id)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { toast } from '@/core/components/ui/toast'
import type { IngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {ingredientCategories} = defineProps<{ingredientCategories: IngredientCategoryDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => ingredientsService.deleteIngredientCategory(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-ingredient-categories'] })
	},
	onError: () => {
		toast({ title: 'Произошла ошибка при удалении' })
	},
})

const onDeleteClick = (e: Event, id: number) => {
	e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить этот продукт?')
	if (confirmed) {
		deleteMutation(id)
	}
}


const goToDetails = (ingredientCategoryId: number) => {
  router.push(`/admin/ingredient-categories/${ingredientCategoryId}`);
};
</script>
