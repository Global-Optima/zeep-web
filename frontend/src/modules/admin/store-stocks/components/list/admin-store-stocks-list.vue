<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Количество</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
				<TableHead
					class="w-12"
					v-if="showForStores"
				></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<!-- If no stocks -->
			<TableRow v-if="stocks.length === 0">
				<TableCell
					colspan="6"
					class="py-6 text-muted-foreground text-center"
				>
					Нет данных
				</TableCell>
			</TableRow>

			<!-- Otherwise, render each stock row -->
			<TableRow
				v-for="stock in stocks"
				:key="stock.id"
				class="hover:bg-gray-50 cursor-pointer"
				@click="handleRowClick(stock.id)"
			>
				<TableCell class="py-4 font-medium">{{ stock.name }}</TableCell>
				<TableCell>{{ stock.ingredient.category.name }}</TableCell>
				<TableCell>{{ stock.quantity }} {{ stock.ingredient.unit.name.toLowerCase() }}</TableCell>

				<!-- Status badge (hidden on small screens) -->
				<TableCell class="hidden md:table-cell">
					<p
						class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
						:class="getStatusClass(stock)"
					>
						{{ getStatusLabel(stock) }}
					</p>
				</TableCell>

				<!-- "Add to cart" button -->
				<TableCell v-if="showForStores">
					<!-- Stop event so row click doesn't fire -->
					<Button
						size="icon"
						variant="ghost"
						aria-label="Добавить сырье в заказ"
						@click.stop="handleAddToCartClick(stock.ingredient.id)"
					>
						<TooltipProvider>
							<Tooltip>
								<TooltipTrigger>
									<PackagePlus class="w-6 h-6 text-gray-500" />
								</TooltipTrigger>
								<TooltipContent>
									<p>Добавить сырье в заказ</p>
								</TooltipContent>
							</Tooltip>
						</TooltipProvider>
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>

	<!-- Dialog for selecting a stock material if needed -->
	<AdminStockMaterialWithQuantitySelectDialog
		:initial-filter="stockMaterialsFilter"
		:open="selectDialogOpen"
		@close="selectDialogOpen = false"
		@submit="handleSelectStockMaterial"
	/>
</template>

<script setup lang="ts">
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { Button } from '@/core/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/core/components/ui/tooltip'

import { PackagePlus } from 'lucide-vue-next'

import type { StockRequestStockMaterialDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/stock-requests/services/stock-requests.service'

import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminStockMaterialWithQuantitySelectDialog from '@/modules/admin/stock-materials/components/admin-stock-material-with-quantity-select-dialog.vue'
import type { StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'

defineProps<{
  stocks: StoreWarehouseStockDTO[]
}>()


const router = useRouter()
function handleRowClick(stockId: number): void {
  router.push(`/admin/store-stocks/${stockId}`)
}

const showForStores = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA])


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

function computeStatus(stock: StoreWarehouseStockDTO): IngredientStatus {
  if (stock.quantity === 0) {
    return 'out_of_stock'
  }
  if (stock.quantity <= stock.lowStockThreshold) {
    return 'low_stock'
  }
  return 'in_stock'
}

function getStatusClass(stock: StoreWarehouseStockDTO): string {
  return INGREDIENT_STATUS_COLOR[computeStatus(stock)]
}

function getStatusLabel(stock: StoreWarehouseStockDTO): string {
  return INGREDIENT_STATUS_FORMATTED[computeStatus(stock)]
}

const queryClient = useQueryClient()

const { mutate: mutateAddMaterial } = useMutation({
  mutationFn: (item: StockRequestStockMaterialDTO) =>
    stockRequestsService.addStockMaterialToLatestCart(item),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
  },
  onError: (error) => {
    console.error('Action failed:', error)
  },
})

const selectDialogOpen = ref(false)
const stockMaterialsFilter = ref<StockMaterialsFilter>({})

function handleAddToCartClick(ingredientId: number): void {
  stockMaterialsFilter.value = { ...stockMaterialsFilter.value, ingredientId }
  selectDialogOpen.value = true
}

function handleSelectStockMaterial(dto: StockRequestStockMaterialDTO): void {
  mutateAddMaterial(dto)
}
</script>

<style scoped></style>
