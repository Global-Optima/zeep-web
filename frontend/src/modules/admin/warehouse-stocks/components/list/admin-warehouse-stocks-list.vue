<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Проект</TableHead>
				<TableHead>Позиция</TableHead>
				<TableHead>Количество</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<!-- If no stocks -->
			<TableRow v-if="stocks.length === 0">
				<TableCell
					colspan="7"
					class="py-6 text-center text-muted-foreground"
				>
					Нет данных
				</TableCell>
			</TableRow>

			<!-- Render each stock row -->
			<TableRow
				v-for="stock in formattedStocks"
				:key="stock.stockMaterial.id"
				class="hover:bg-gray-50 cursor-pointer"
				@click="handleRowClick(stock.stockMaterial.id)"
			>
				<TableCell class="py-4 font-medium">{{ stock.stockMaterial.name }}</TableCell>
				<TableCell>{{ stock.projectName }}</TableCell>
				<TableCell>{{ stock.position }}</TableCell>
				<TableCell>{{ stock.quantity }}</TableCell>

				<!-- Status badges -->
				<TableCell class="hidden md:table-cell">
					<div class="flex flex-col space-y-2">
						<!-- Stock Status -->
						<p
							class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
							:class="getStatusClass(stock)"
						>
							{{ getStatusLabel(stock) }}
						</p>

						<!-- Expiration Indicator -->
						<p
							v-if="shouldIndicateExpiration(stock)"
							class="inline-flex items-center bg-orange-100 px-2.5 py-1 rounded-md w-fit text-orange-800 text-xs"
						>
							Истекает через 7 дней
						</p>
					</div>
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
import type { WarehouseStocksDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { differenceInDays } from 'date-fns'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

interface TableType extends WarehouseStocksDTO {
  projectName: string, position: string
}

// Define props
const {stocks} = defineProps<{
	stocks: WarehouseStocksDTO[]
}>()

// Mock projects and positions for the stocks
const mockProjects = ['Модернизация системы учета', 'Поставка оборудования для склада', 'Ремонта оборудования', 'Обучение персонала']
const mockPositions = ['A8', 'B7', 'C3', 'D4', 'E5', 'E8', 'B3', 'H7']

// Enhance stocks with projects and positions
const formattedStocks = computed<TableType[]>(() => stocks.map((stock, index) => ({
  ...stock,
	projectName: mockProjects[index % mockProjects.length],
	position: mockPositions[index % mockPositions.length]
})))

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
	if (stock.quantity === 0) {
		return 'out_of_stock'
	}
	if (stock.quantity <= stock.stockMaterial.safetyStock) {
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

function shouldIndicateExpiration(stock: WarehouseStocksDTO): boolean {
	if (!stock.earliestExpirationDate) {
		return false
	}
	const daysLeft = differenceInDays(new Date(stock.earliestExpirationDate), new Date())
	return daysLeft <= 7 && daysLeft > 0
}
</script>
