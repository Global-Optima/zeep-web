<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Обьем</TableHead>
				<TableHead>Себестоимость</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="provision in provisions"
				:key="provision.id"
				@click="onProvisionClick(provision.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ provision.name }}</TableCell>
				<TableCell
					>{{ provision.absoluteVolume }} {{ provision.unit.name.toLowerCase() }}</TableCell
				>
				<TableCell>{{ provision.netCost }}</TableCell>
				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, provision.id)"
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
import type { ProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {provisions} = defineProps<{provisions: ProvisionDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => provisionsService.deleteProvision(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-provisions'] })
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

const onProvisionClick = (id: number) => {
  router.push(`/admin/provisions/${id}`);
};
</script>
