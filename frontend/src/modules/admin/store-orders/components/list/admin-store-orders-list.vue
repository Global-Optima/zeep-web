<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Создано</TableHead>
				<TableHead class="p-4">Заказчик</TableHead>
				<TableHead class="p-4">Сумма</TableHead>
				<!-- <TableHead class="hidden md:table-cell p-4">Оплата</TableHead> -->
				<TableHead class="hidden md:table-cell p-4">Статус</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="order in orders"
				:key="order.id"
				class="hover:bg-gray-50 h-12"
			>
				<TableCell class="p-4">
					{{ formatDate(order.createdAt) }}
				</TableCell>
				<TableCell class="p-4 font-medium">
					{{ order.customerName }}
				</TableCell>
				<TableCell class="p-4 font-medium">
					{{ formatPrice(order.total) }}
				</TableCell>

				<!-- TODO: add paid cell back -->
				<!-- <TableCell class="hidden md:table-cell p-4">
					<div
						v-if="order?.isPaid"
						class="flex items-center gap-2 text-green-500"
					>
						<Check class="w-4 h-4" />
						Оплачено
					</div>
					<div
						v-else
						class="flex items-center gap-2 text-red-500"
					>
						<X class="w-4 h-4" />
						Не оплачено
					</div>
				</TableCell> -->

				<TableCell class="hidden md:table-cell p-4">
					<p
						:class="[
                'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
                ORDER_STATUS_COLOR[order.status],
              ]"
					>
						{{ ORDER_STATUS_FORMATTED[order.status] }}
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
import { OrderStatus, type OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'

// Props
const { orders } = defineProps<{ orders: OrderDTO[] }>();

const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'dd MMMM yyyy', { locale: ru })
}
const ORDER_STATUS_COLOR: Record<OrderStatus, string> = {
  [OrderStatus.PENDING]: 'bg-purple-100 text-purple-800',
  [OrderStatus.PREPARING]: 'bg-yellow-100 text-yellow-800',
  [OrderStatus.COMPLETED]: 'bg-green-200 text-green-900',
  [OrderStatus.IN_DELIVERY]: 'bg-orange-100 text-orange-800',
  [OrderStatus.DELIVERED]: 'bg-green-100 text-green-800',
  [OrderStatus.CANCELLED]: 'bg-red-100 text-red-800'
}

const ORDER_STATUS_FORMATTED: Record<OrderStatus, string> = {
  [OrderStatus.PENDING]: 'В ожидании',
  [OrderStatus.PREPARING]: 'Готовится',
  [OrderStatus.COMPLETED]: 'Завершен',
  [OrderStatus.IN_DELIVERY]: 'Доставляется',
  [OrderStatus.DELIVERED]: 'Доставлен',
  [OrderStatus.CANCELLED]: 'Отменен'
}
</script>
