<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Сырье</TableHead>
				<TableHead>Упаковка</TableHead>
				<TableHead
					v-if="canDelete"
					class="p-4"
				></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="stockMaterial in stockMaterials"
				:key="stockMaterial.id"
				@click="onStockMaterialClick(stockMaterial.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ stockMaterial.name }}</TableCell>
				<TableCell>{{ stockMaterial.category.name }}</TableCell>
				<TableCell>{{ stockMaterial.ingredient.name }}</TableCell>
				<TableCell>{{ stockMaterial.size }} {{ stockMaterial.unit.name.toLowerCase() }}</TableCell>
				<TableCell class="flex justify-end">
					<Button
						v-if="canDelete"
						variant="ghost"
						size="icon"
						@click="(e: Event) => onDeleteClick(e, stockMaterial.id)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
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
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { useRouter } from 'vue-router'
import { stockMaterialsService } from '../../services/stock-materials.service';
import { useMutation, useQueryClient } from '@tanstack/vue-query';
import { toast } from '@/core/components/ui/toast';
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks';
import { useHasRole } from '@/core/hooks/use-has-roles.hook';
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models';
import { Trash } from 'lucide-vue-next'
import { Button } from '@/core/components/ui/button'


const {stockMaterials} = defineProps<{stockMaterials: StockMaterialsDTO[]}>()
const { toastLocalizedError } = useAxiosLocaleToast()

const router = useRouter();
const queryClient = useQueryClient()

const canDelete = useHasRole([EmployeeRole.ADMIN])

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => stockMaterialsService.deleteStockMaterial(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-stock-materials'] })
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

const onStockMaterialClick = (id: number) => {
  router.push(`/admin/stock-materials/${id}`);
};
</script>
