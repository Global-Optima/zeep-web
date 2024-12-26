<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Цена</TableHead>
				<TableHead class="hidden md:table-cell">Доступно</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="product in products"
				:key="product.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
			>
				<TableCell class="font-medium">
					{{ product.name }}
				</TableCell>
				<TableCell>
					{{ product.category }}
				</TableCell>
				<TableCell>
					{{ formatPrice(product.price) }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					<div
						:class="[
								product.isAvailable ? 'text-green-500' : 'text-red-500',
							]"
					>
						{{ product.isAvailable ? 'В наличии' : 'Нет в наличии' }}
					</div>
				</TableCell>
				<TableCell class="hidden md:table-cell">
					<p
						:class="[
								'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
								PRODUCT_STATUS_COLOR[product.status]
							]"
					>
						{{ PRODUCT_STATUS_FORMATTED[product.status] }}
					</p>
				</TableCell>
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
import { ref } from 'vue'

const products = ref([
    {
      id: 1,
      name: 'Латте',
      category: 'Кофе',
      price: 150.0,
      isAvailable: true,
      status: 'active',
    },
    {
      id: 2,
      name: 'Американо',
      category: 'Кофе',
      price: 100.0,
      isAvailable: false,
      status: 'inactive',
    },
    {
      id: 3,
      name: 'Чай зеленый',
      category: 'Чай',
      price: 80.0,
      isAvailable: true,
      status: 'active',
    },
    {
      id: 4,
      name: 'Смузи клубничный',
      category: 'Напитки',
      price: 250.0,
      isAvailable: true,
      status: 'active',
    },
    {
      id: 5,
      name: 'Круассан с сыром',
      category: 'Выпечка',
      price: 120.0,
      isAvailable: false,
      status: 'inactive',
    },
    {
      id: 6,
      name: 'Молочный коктейль',
      category: 'Напитки',
      price: 200.0,
      isAvailable: true,
      status: 'active',
    },
])

const PRODUCT_STATUS_COLOR: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-gray-100 text-gray-800',
}

const PRODUCT_STATUS_FORMATTED: Record<string, string> = {
  active: 'Активен',
  inactive: 'Неактивен',
}
</script>

<style scoped></style>
