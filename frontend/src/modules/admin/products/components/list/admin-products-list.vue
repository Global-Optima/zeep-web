<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden w-[100px] sm:table-cell"> </TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Начальная цена</TableHead>
				<TableHead>Размеры</TableHead>
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
						class="bg-gray-100 rounded-md aspect-square object-contain"
						height="64"
						width="64"
					/>
				</TableCell>
				<TableCell class="font-medium">{{ product.name }}</TableCell>
				<TableCell>{{ product.categoryName }}</TableCell>
				<TableCell>{{ formatPrice(product.basePrice) }}</TableCell>
				<TableCell>{{ product.productSizeCount }}</TableCell>
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
import { formatPrice } from '@/core/utils/price.utils'
import type { ProductDTO } from '@/modules/kiosk/products/models/product.model'
import { useRouter } from 'vue-router'

const {products} = defineProps<{products: ProductDTO[]}>()

const router = useRouter();

const onProductClick = (productId: number) => {
  router.push(`/admin/products/${productId}`);
};
</script>
