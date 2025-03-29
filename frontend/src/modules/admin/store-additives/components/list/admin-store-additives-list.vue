<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden sm:table-cell w-[100px]"> </TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Размер</TableHead>
				<TableHead>Цена</TableHead>
        <TableHead>Статус</TableHead>
				<TableHead v-if="canDelete"></TableHead>
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
					<LazyImage
						:src="additive.imageUrl"
						alt="Изображение добавки"
						class="rounded-md size-16 object-contain aspect-square"
					/>
				</TableCell>
				<TableCell class="font-medium">
					{{ additive.name }}
				</TableCell>
				<TableCell>
					{{ additive.category.name }}
				</TableCell>
				<TableCell> {{ additive.size }} {{ additive.unit.name.toLowerCase() }} </TableCell>
				<TableCell>
					{{ formatPrice(additive.storePrice) }}
				</TableCell>
        <TableCell class="hidden md:table-cell">
          <p
            class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
            :class="getStatusClass(additive)"
          >
            {{ getStatusLabel(additive) }}
          </p>
        </TableCell>
				<TableCell v-if="canDelete">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteProductClick(e, additive.id)"
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
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { formatPrice } from '@/core/utils/price.utils'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { StoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

type IngredientStatus = 'in_stock' | 'out_of_stock'

const STORE_ADDITIVE_STATUS_COLOR: Record<IngredientStatus, string> = {
  in_stock: 'bg-green-100 text-green-800',
  out_of_stock: 'bg-red-100 text-red-800',
}

const STORE_ADDITIVE_STATUS_FORMATTED: Record<IngredientStatus, string> = {
  in_stock: 'В наличии',
  out_of_stock: 'Нет в наличии',
}

function computeStatus(additive: StoreAdditiveDTO): IngredientStatus {
  if (additive.isOutOfStock) {
    return 'out_of_stock'
  }
  return 'in_stock'
}

function getStatusClass(additive: StoreAdditiveDTO): string {
  return STORE_ADDITIVE_STATUS_COLOR[computeStatus(additive)]
}

function getStatusLabel(additive: StoreAdditiveDTO): string {
  return STORE_ADDITIVE_STATUS_FORMATTED[computeStatus(additive)]
}

const router = useRouter()
const queryClient = useQueryClient()

const {additives} = defineProps<{additives: StoreAdditiveDTO[]}>()

const canDelete = useHasRole([EmployeeRole.STORE_MANAGER])

const {mutate: deleteStoreAdditive} = useMutation({
		mutationFn: (id: number) => storeAdditivesService.deleteStoreAdditive(id),
		onSuccess: () => {
			toast({title: "Топпинг удален из кафе"})
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
  if (!canDelete) return

  e.stopPropagation()
  deleteStoreAdditive(id)
}
</script>

<style scoped></style>
