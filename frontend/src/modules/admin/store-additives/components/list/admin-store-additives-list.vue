<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead class="hidden md:table-cell">Описание</TableHead>
				<TableHead>Цена</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="additive in additives"
				:key="additive.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="onProductClick(additive.storeAdditiveId)"
			>
				<TableCell class="font-medium">
					{{ additive.name }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					{{ additive.description }}
				</TableCell>
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
