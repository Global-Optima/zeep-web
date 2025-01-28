<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden w-[100px] sm:table-cell"> </TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Размер</TableHead>
				<TableHead>Цена</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="additive in additives"
				:key="additive.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="onProductClick(additive.id)"
			>
				<TableCell class="hidden sm:table-cell">
					<img
						:src="additive.imageUrl"
						alt="Изображение товара"
						class="bg-gray-100 p-1 rounded-md aspect-square object-contain"
						height="64"
						width="64"
					/>
				</TableCell>
				<TableCell class="font-medium">
					{{ additive.name }}
				</TableCell>
				<TableCell>
					{{ additive.category.name }}
				</TableCell>
				<TableCell> {{ additive.size }} {{ additive.unit.name }} </TableCell>
				<TableCell>
					{{ formatPrice(additive.storePrice) }}
				</TableCell>
				<TableCell>
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteProductClick(e, additive.storeAdditiveId)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import Button from '@/core/components/ui/button/Button.vue'
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
import type { StoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const {additives} = defineProps<{additives: StoreAdditiveDTO[]}>()

const {mutate: deleteStoreAdditive} = useMutation({
		mutationFn: (id: number) => storeAdditivesService.deleteStoreAdditive(id),
		onSuccess: () => {
			toast({title: "Топпинг удален из магазина"})
			queryClient.invalidateQueries({queryKey: ['admin-store-additives']})
		},
		onError: () => {
      toast({title: "Произошла ошибка при удалении"})
		},
})

const onProductClick = (additiveId: number) => {
  router.push(`/admin/store-additives/${additiveId}`);
};

const onDeleteProductClick = (e: Event, id: number) => {
  e.stopPropagation()
  deleteStoreAdditive(id)
}
</script>

<style scoped></style>
