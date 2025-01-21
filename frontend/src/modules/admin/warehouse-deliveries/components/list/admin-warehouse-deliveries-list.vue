<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Поставщик</TableHead>
				<TableHead>Дата Доставки</TableHead>
				<TableHead>Товаров</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="delivery in deliveries"
				:key="delivery.id"
				class="hover:bg-gray-50 cursor-pointer"
				@click="handleRowClick(delivery.id)"
			>
				<TableCell class="py-4 font-medium">{{ delivery.supplier.name }}</TableCell>
				<TableCell>{{ format(delivery.deliveryDate, "DD.MM.YYYY") }}</TableCell>
				<TableCell>{{  delivery.materials.length }}</TableCell>
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
import type { WarehouseDeliveriesDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { format } from 'date-fns'


defineProps<{
  deliveries: WarehouseDeliveriesDTO[]
}>()


const router = useRouter()
function handleRowClick(stockId: number): void {
  router.push(`/admin/warehouse-deliveries/${stockId}`)
}
</script>

<style scoped></style>
