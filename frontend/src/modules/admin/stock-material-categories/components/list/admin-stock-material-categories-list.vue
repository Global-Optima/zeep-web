<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead
					v-if="canDelete"
					class="p-4"
				></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="category in categories"
				:key="category.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(category.id)"
			>
				<TableCell class="p-4">
					<span class="font-medium">{{ category.name }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span class="line-clamp-1 max-w-md">{{ category.description }}</span>
				</TableCell>

				<TableCell class="flex justify-end">
					<Button
						v-if="canDelete"
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
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
const { categories} = defineProps<{categories: StockMaterialCategoryDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const canDelete = useHasRole([EmployeeRole.ADMIN])
const { toastLocalizedError } = useAxiosLocaleToast()


const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => stockMaterialCategoryService.delete(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories'] })
	},
	onError: (error: AxiosLocalizedError) => {
		toastLocalizedError(error, 'Произошла ошибка при удалении' )
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
  router.push(`/admin/stock-material-categories/${productCategoryId}`);
};
</script>
