<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden w-[100px] sm:table-cell"></TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Начальная цена</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="product in products"
				:key="product.id"
				@click="onProductClick(product.id)"
				class="hover:bg-slate-50 cursor-pointer"
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
				<TableCell class="font-medium">{{ product.name }}</TableCell>
				<TableCell>{{ product.category.name }}</TableCell>
				<TableCell>{{ formatPrice(product.basePrice) }}</TableCell>
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
import { toast } from '@/core/components/ui/toast'
import { formatPrice } from '@/core/utils/price.utils'
import type { ProductDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const { products } = defineProps<{ products: ProductDTO[] }>()

const router = useRouter()
const queryClient = useQueryClient()

const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => productsService.deleteProduct(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-products'] })
	},
	onError: () => {
		toast({ title: 'Произошла ошибка при удалении' })
	},
})

const onDeleteClick = (e: Event, id: number) => {
	e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить этот продукт?')
	if (confirmed) {
		deleteMutation(id)
	}
}

const onProductClick = (productId: number) => {
	router.push(`/admin/products/${productId}`)
}
</script>
