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
			<!-- If no stocks -->
			<TableRow v-if="stocks.length === 0">
				<TableCell
					colspan="6"
					class="py-6 text-center text-muted-foreground"
				>
					Нет данных
				</TableCell>
			</TableRow>

			<!-- Otherwise, render each stock row -->
			<TableRow
				v-for="stock in stocks"
				:key="stock.stockMaterial.id"
				class="hover:bg-gray-50 cursor-pointer"
				@click="handleRowClick(stock.stockMaterial.id)"
			>
				<TableCell class="py-4 font-medium">{{ stock.stockMaterial.name }}</TableCell>
				<TableCell>{{ stock.totalQuantity }}</TableCell>
				<TableCell>{{ stock.stockMaterial.safetyStock }}</TableCell>

				<!-- Unit name (hidden on small screens) -->
				<TableCell class="hidden md:table-cell">
					<!-- {{ stock.stockMaterial.ingredient.unit.name }} -->
				</TableCell>

				<!-- Status badge (hidden on small screens) -->
				<TableCell class="hidden md:table-cell">
					<p
						class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
						:class="getStatusClass(stock)"
					>
						{{ getStatusLabel(stock) }}
					</p>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'
import type { WarehouseStocksDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'


defineProps<{
  stocks: WarehouseStocksDTO[]
}>()


const router = useRouter()
function handleRowClick(stockId: number): void {
  router.push(`/admin/warehouse-stocks/${stockId}`)
}

type IngredientStatus = 'in_stock' | 'low_stock' | 'out_of_stock'

const INGREDIENT_STATUS_COLOR: Record<IngredientStatus, string> = {
  in_stock: 'bg-green-100 text-green-800',
  low_stock: 'bg-yellow-100 text-yellow-800',
  out_of_stock: 'bg-red-100 text-red-800',
}

const INGREDIENT_STATUS_FORMATTED: Record<IngredientStatus, string> = {
  in_stock: 'В наличии',
  low_stock: 'Заканчивается',
  out_of_stock: 'Нет в наличии',
}


function computeStatus(stock: WarehouseStocksDTO): IngredientStatus {
  if (stock.totalQuantity === 0) {
    return 'out_of_stock'
  }
  if (stock.totalQuantity <= stock.stockMaterial.safetyStock) {
    return 'low_stock'
  }
  return 'in_stock'
}

function getStatusClass(stock: WarehouseStocksDTO): string {
  return INGREDIENT_STATUS_COLOR[computeStatus(stock)]
}

function getStatusLabel(stock: WarehouseStocksDTO): string {
  return INGREDIENT_STATUS_FORMATTED[computeStatus(stock)]
}
</script>

<style scoped></style>
