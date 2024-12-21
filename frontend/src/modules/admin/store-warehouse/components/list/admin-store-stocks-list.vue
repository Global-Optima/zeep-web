<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Количество</TableHead>
				<TableHead>Мин. запас</TableHead>
				<TableHead class="hidden md:table-cell">Единица измерения</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="stock in stocks"
				:key="stock.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="goToDetails(stock.id)"
			>
				<TableCell class="font-medium">
					{{ stock.name }}
				</TableCell>
				<TableCell>
					{{ stock.quantity }}
				</TableCell>
				<TableCell>
					{{ stock.lowStockThreshold }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					{{ stock.unit }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					<p
						:class="[
								'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
								INGREDIENT_STATUS_COLOR[getStockStatus(stock)]
							]"
					>
						{{ INGREDIENT_STATUS_FORMATTED[getStockStatus(stock)] }}
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
import type { StoreStocks } from '@/modules/admin/store-warehouse/models/store-stock.model'
import { useRouter } from 'vue-router'

// Props
const { stocks } = defineProps<{ stocks: StoreStocks[] }>();

// Router for details navigation
const router = useRouter();

const goToDetails = (stockId: number) => {
  router.push(`/admin/store-stocks/${stockId}`);
};

// Status mapping
const INGREDIENT_STATUS_COLOR: Record<string, string> = {
  in_stock: 'bg-green-100 text-green-800',
  low_stock: 'bg-yellow-100 text-yellow-800',
  out_of_stock: 'bg-red-100 text-red-800',
};

const INGREDIENT_STATUS_FORMATTED: Record<string, string> = {
  in_stock: 'В наличии',
  low_stock: 'Заканчивается',
  out_of_stock: 'Нет в наличии',
};

// Function to determine the stock status
const getStockStatus = (stock: StoreStocks): string => {
  if (stock.quantity === 0) {
    return 'out_of_stock';
  }
  if (stock.quantity <= stock.lowStockThreshold) {
    return 'low_stock';
  }
  return 'in_stock';
};
</script>
