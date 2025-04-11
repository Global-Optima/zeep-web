<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Множитель</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="unit in units"
				:key="unit.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(unit.id)"
			>
				<TableCell class="p-4 font-medium"> {{ unit.name }} </TableCell>
				<TableCell> {{ unit.conversionFactor }} </TableCell>
				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, unit.id)"
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
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {units} = defineProps<{units: UnitDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()
const { toastLocalizedError } = useAxiosLocaleToast()


const {mutate: deleteMutation} = useMutation({
		mutationFn: (id: number) => unitsService.deleteUnit(id),
		onSuccess: () => {
			toast({title: "Успешное удаление"})
			queryClient.invalidateQueries({queryKey: ['admin-units']})
		},
		onError: (error: AxiosLocalizedError) => {
			toastLocalizedError(error, "Произошла ошибка при удалении")
		},
})

const onDeleteClick = (e: Event, id: number) => {
  e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить?')
	if (confirmed) {
		deleteMutation(id)
	}}

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/units/${productCategoryId}`);
};
</script>
