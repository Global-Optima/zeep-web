<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden sm:table-cell w-[100px]"></TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Базовая цена</TableHead>
				<TableHead>Размер</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="additive in additives"
				:key="additive.id"
				@click="onAdditiveClick(additive.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="hidden sm:table-cell">
					<LazyImage
						:src="additive.imageUrl"
						alt="Изображение добавки"
						class="rounded-md size-16 object-contain aspect-square"
					/>
				</TableCell>
				<TableCell class="font-medium">{{ additive.name }}</TableCell>
				<TableCell>{{ additive.category.name }}</TableCell>
				<TableCell>{{ formatPrice(additive.basePrice) }}</TableCell>
				<TableCell>{{ additive.size }}</TableCell>
				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, additive.id)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
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
import { formatPrice } from '@/core/utils/price.utils'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {additives} = defineProps<{additives: AdditiveDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => additivesService.deleteAdditive(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
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

const onAdditiveClick = (additiveID: number) => {
  router.push(`/admin/additives/${additiveID}`);
};
</script>
