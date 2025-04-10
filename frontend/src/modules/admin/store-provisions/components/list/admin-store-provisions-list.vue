<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Обьем</TableHead>
				<TableHead>Создан</TableHead>
				<TableHead>Статус</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="storeProvision in storeProvisions"
				:key="storeProvision.id"
				@click="onProvisionClick(storeProvision.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ storeProvision.provision.name }}</TableCell>
				<TableCell
					>{{ storeProvision.volume }}
					{{ storeProvision.provision.unit.name.toLowerCase() }}</TableCell
				>
				<TableCell>{{ formatCreatedDate(storeProvision.createdAt) }}</TableCell>

				<TableCell>
					<p
						class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
						:class="STORE_PROVISION_STATUS_COLOR[storeProvision.status]"
					>
						{{ STORE_PROVISION_STATUS_FORMATTED[storeProvision.status] }}
					</p>
				</TableCell>

				<TableCell class="flex justify-end">
					<Button
						v-if="canDelete(storeProvision)"
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, storeProvision.id)"
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
import { StoreProvisionStatus, type StoreProvisionDTO } from "@/modules/admin/store-provisions/models/store-provision.models"
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {storeProvisions} = defineProps<{storeProvisions: StoreProvisionDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const {toastLocalizedError} = useAxiosLocaleToast()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => storeProvisionsService.deleteStoreProvision(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-store-provisions'] })
	},
	onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, 'Произошла ошибка при удалении')
	},
})

const onDeleteClick = (e: Event, id: number) => {
	e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить?')
	if (confirmed) {
		deleteMutation(id)
	}
}

const formatCreatedDate = (date: Date) => {
  return format(date, "dd.MM.yyyy hh:mm", {locale: ru})
}

const canDelete = (provision: StoreProvisionDTO) => {
  return provision.status === StoreProvisionStatus.PREPARING
}

const onProvisionClick = (id: number) => {
  router.push(`/admin/store-provisions/${id}`);
};

const STORE_PROVISION_STATUS_COLOR: Record<StoreProvisionStatus, string> = {
  [StoreProvisionStatus.COMPLETED]: 'bg-green-100 text-green-800',
  [StoreProvisionStatus.PREPARING]: 'bg-yellow-100 text-yellow-800',
  [StoreProvisionStatus.EXPIRED]: 'bg-red-100 text-red-800',
}

const STORE_PROVISION_STATUS_FORMATTED: Record<StoreProvisionStatus, string> = {
  [StoreProvisionStatus.COMPLETED]: 'Приготовлен',
  [StoreProvisionStatus.PREPARING]: 'Готовится',
  [StoreProvisionStatus.EXPIRED]: 'Просрочен',
}
</script>
