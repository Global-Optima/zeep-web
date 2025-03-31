<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead class="p-4">Кол-во топпингов</TableHead>
				<TableHead class="p-4"></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="category in categories"
				:key="category.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(category.id)"
			>
				<TableCell class="p-4 font-medium"> {{ category.name }} </TableCell>
				<TableCell> {{ category.description }} </TableCell>
				<TableCell> {{ category.additivesCount }} </TableCell>
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
import type { AdditiveCategoryDetailsDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {categories} = defineProps<{categories: AdditiveCategoryDetailsDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => additivesService.deleteAdditiveCategory(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories'] })
	},
	onError: () => {
		toast({ title: 'Произошла ошибка при удалении' })
	},
})

const onDeleteClick = (e: Event, id: number) => {
	e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить?')
	if (confirmed) {
		deleteMutation(id)
	}
}

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/additive-categories/${productCategoryId}`);
};
</script>
