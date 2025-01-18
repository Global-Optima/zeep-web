<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Количество</TableHead>
				<TableHead>Мин. запас</TableHead>
				<TableHead class="hidden md:table-cell">Единица измерения</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
				<TableHead></TableHead>
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
					{{ stock.ingredient.unit.name }}
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
				<TableCell>
					<Button
						size="icon"
						variant="ghost"
						@click="e => onAddToCartClick(e, stock.ingredient.id)"
					>
						<TooltipProvider>
							<Tooltip>
								<TooltipTrigger> <PackagePlus class="w-6 h-6 text-gray-500" /> </TooltipTrigger>
								<TooltipContent>
									<p>Добавить ингредиент в заказ</p>
								</TooltipContent>
							</Tooltip>
						</TooltipProvider>
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>

	<AdminStockMaterialsSelectDialog
		:initial-filter="stockMaterialsFilter"
		:open="openStockMaterialsDialog"
		@close="openStockMaterialsDialog = false"
		@select="onSelectStockMaterial"
	/>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/core/components/ui/tooltip'
import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'
import type { StockMaterialsDTO, StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import { PackagePlus } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

// Props
const { stocks } = defineProps<{ stocks: StoreWarehouseStockDTO[] }>();

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
const getStockStatus = (stock: StoreWarehouseStockDTO): string => {
  if (stock.quantity === 0) {
    return 'out_of_stock';
  }
  if (stock.quantity <= stock.lowStockThreshold) {
    return 'low_stock';
  }
  return 'in_stock';
};

const stockMaterialsFilter = ref<StockMaterialsFilter>({})
const openStockMaterialsDialog = ref(false)

const onAddToCartClick = (e: Event, ingredientId: number) => {
  e.stopPropagation();
  stockMaterialsFilter.value = { ...stockMaterialsFilter.value, ingredientId };
  openStockMaterialsDialog.value = true;
};

const onSelectStockMaterial = (stockMaterial: StockMaterialsDTO) => {
  console.log(stockMaterial)
}
</script>
