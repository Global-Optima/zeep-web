<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Поставщик</TableHead>
				<TableHead>Дата доставки</TableHead>
				<TableHead>Продуктов</TableHead>
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
				<TableCell>{{ format(delivery.deliveryDate, "d MMMM yyyy", {locale:ru}) }}</TableCell>
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
import type { WarehouseDeliveryDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'


defineProps<{
  deliveries: WarehouseDeliveryDTO[]
}>()


const router = useRouter()
function handleRowClick(stockId: number): void {
  router.push(`/admin/warehouse-deliveries/${stockId}`)
}
</script>

<style scoped></style>
