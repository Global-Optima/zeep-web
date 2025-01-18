<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Количество</TableHead>
				<!-- <TableHead>Мин. запас</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead> -->
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="stock in stocks"
				:key="stock.stockMaterialId"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="goToDetails(stock.stockMaterialId)"
			>
				<TableCell class="font-medium">
					{{ stock.name }}
				</TableCell>
				<TableCell>
					{{ stock.quantity }}
				</TableCell>
				<!-- <TableCell>
					{{ stock.safetyStock }}
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
				</TableCell> -->
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
import type { InventoryLevel } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { useRouter } from 'vue-router'

// Props
const { stocks } = defineProps<{ stocks: InventoryLevel[] }>();

// Router for details navigation
const router = useRouter();

const goToDetails = (stockId: number) => {
  router.push(`/admin/warehouse-stocks/${stockId}`);
};
</script>

<style scoped></style>
