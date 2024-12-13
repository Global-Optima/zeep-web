<template>
	<main class="flex-1 items-start gap-4 grid">
		<div class="flex justify-between items-center gap-2">
			<div class="flex items-center gap-2">
				<Input
					placeholder="Поиск товаров"
					class="bg-white w-full md:w-64"
				/>
				<DropdownMenu>
					<DropdownMenuTrigger as-child>
						<Button
							variant="outline"
							class="gap-2"
						>
							<ListFilter class="w-4 h-4" />
							<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Фильтр</span>
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuLabel>Фильтровать по</DropdownMenuLabel>
						<DropdownMenuSeparator />
						<DropdownMenuItem checked>Активные</DropdownMenuItem>
						<DropdownMenuItem>Черновики</DropdownMenuItem>
						<DropdownMenuItem>Архив</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</div>

			<div class="flex items-center gap-2">
				<Button variant="outline">
					<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Экспорт</span>
				</Button>
				<Button>
					<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Добавить</span>
				</Button>
			</div>
		</div>

		<Card>
			<CardContent class="mt-4">
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
							<TableCell>{{ product.category }}</TableCell>
							<TableCell>{{ formatPrice(product.startPrice) }}</TableCell>
							<TableCell>{{ product.sizesQuantity }}</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</main>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent
} from '@/core/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { getRoutePath } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { ListFilter } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const onProductClick = (productId: number) => {
  router.push(`${getRoutePath("ADMIN_PRODUCTS")}/${productId}`)
}

// Mock data for products
const products = ref([
  {
    id: 1,
    imageUrl: 'fake_url',
    name: 'Латте',
    category: 'Кофе',
    startPrice: 1200,
    sizesQuantity: 3,
  },
  {
    id: 2,
    imageUrl: 'fake_url',
    name: 'Американо',
    category: 'Кофе',
    startPrice: 1100,
    sizesQuantity: 3,
  },
  {
    id: 3,
    imageUrl: 'fake_url',
    name: 'Чай зеленый',
    category: 'Чай',
    startPrice: 850,
    sizesQuantity: 3,
  },
  {
    id: 4,
    imageUrl: 'fake_url',
    name: 'Смузи клубничный',
    category: 'Напитки',
    startPrice: 2000,
    sizesQuantity: 3,
  },
  {
    id: 5,
    imageUrl: 'fake_url',
    name: 'Круассан с сыром',
    category: 'Выпечка',
    startPrice: 2250,
    sizesQuantity: 3,
  },
  {
    id: 6,
    imageUrl: 'fake_url',
    name: 'Молочный коктейль',
    category: 'Напитки',
    startPrice: 1540,
    sizesQuantity: 3,
  },
])
</script>

<style scoped>
/* Add any custom styles here */
</style>
