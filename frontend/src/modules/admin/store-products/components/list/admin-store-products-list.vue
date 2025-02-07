<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden w-[100px] sm:table-cell"> </TableHead>
				<TableHead>Название</TableHead>
				<TableHead class="hidden md:table-cell">Категория</TableHead>
				<TableHead>Цена</TableHead>
				<TableHead class="hidden md:table-cell">Доступно</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="product in storeProducts"
				:key="product.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="onProductClick(product.id)"
			>
				<TableCell class="hidden sm:table-cell">
					<img
						:src="product.imageUrl"
						alt="Изображение товара"
						class="bg-gray-100 p-1 rounded-md aspect-square object-contain"
						height="64"
						width="64"
					/>
				</TableCell>
				<TableCell class="font-medium">
					{{ product.name }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					{{ product.category.name }}
				</TableCell>
				<TableCell>
					{{ formatPrice(product.storePrice) }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					<div
						:class="[
								product.isAvailable ? 'text-green-600' : 'text-red-500',
							]"
					>
						{{ product.isAvailable ? 'Доступен' : 'Недоступен' }}
					</div>
				</TableCell>
				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteProductClick(e, product.id)"
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
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const {storeProducts} = defineProps<{storeProducts: StoreProductDTO[]}>()

const {mutate: deleteStoreProduct} = useMutation({
		mutationFn: (id: number) => storeProductsService.deleteStoreProduct(id),
		onSuccess: () => {
			toast({title: "Товар удален из кафе"})
			queryClient.invalidateQueries({queryKey: ['admin-store-products']})
		},
		onError: () => {
      toast({title: "Произошла ошибка при удалении товара"})
		},
})

const onProductClick = (storeProductId: number) => {
  router.push(`/admin/store-products/${storeProductId}`);
};

const onDeleteProductClick = (e: Event, id: number) => {
  e.stopPropagation()
  deleteStoreProduct(id)
}
</script>

<style scoped></style>
