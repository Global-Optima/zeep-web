<!-- AdminOrdersList.vue -->
<template>
	<div>
		<!-- If no orders, display message -->
		<p
			v-if="orders.data.length === 0"
			class="text-muted-foreground"
		>
			Заказы не найдены
		</p>
		<!-- If there are orders, display the table -->
		<Table v-else>
			<TableHeader>
				<TableRow>
					<TableHead>Создано</TableHead>
					<TableHead>Заказчик</TableHead>
					<TableHead>Сумма</TableHead>
					<TableHead class="hidden md:table-cell">Оплата</TableHead>
					<TableHead class="hidden md:table-cell">Статус</TableHead>
					<TableHead class="hidden md:table-cell">Комментарий</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<TableRow
					v-for="order in displayedOrders"
					:key="order.id"
					class="h-12 cursor-pointer"
					@click="goToOrder(order.id)"
				>
					<TableCell class="font-medium">
						{{ formatDate(order.createdAt) }}
					</TableCell>
					<TableCell class="font-medium">
						{{ getUserFullName(order.user) }}
					</TableCell>
					<TableCell class="font-medium">
						{{ formatPrice(order.totalPrice) }}
					</TableCell>
					<TableCell class="hidden font-medium md:table-cell">
						<div
							v-if="order.isPaid"
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
					</TableCell>
					<TableCell class="hidden font-medium md:table-cell">
						<p
							:class="[
                'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
                ORDER_STATUS_COLOR[order.status],
              ]"
						>
							{{ ORDER_STATUS_FORMATTED[order.status] }}
						</p>
					</TableCell>
					<TableCell class="hidden font-medium md:table-cell">
						<p class="max-w-60 truncate">
							{{ order.notes || 'Нет комментариев' }}
						</p>
					</TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { Check, X } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { formatPrice } from '@/core/utils/price.utils'

// Constants
const DEFAULT_LIMIT = 10

// Mock data for orders (in Russian)
const orders = ref({
  data: [
    {
      id: 1,
      createdAt: '2023-10-01T12:34:56Z',
      user: {
        firstName: 'Иван',
        lastName: 'Иванов',
      },
      totalPrice: 1234.56,
      isPaid: true,
      status: 'completed',
      notes: 'Доставить после 18:00',
    },
    {
      id: 2,
      createdAt: '2023-10-02T08:15:30Z',
      user: {
        firstName: 'Мария',
        lastName: 'Петрова',
      },
      totalPrice: 789.0,
      isPaid: false,
      status: 'pending',
      notes: '',
    },
    // Add more mock orders as needed
  ],
  meta: {
    totalItems: 20, // Total number of orders (for pagination)
  },
})

const limit = ref(DEFAULT_LIMIT)

// Computed property for displayed orders based on limit
const displayedOrders = computed(() => orders.value.data.slice(0, limit.value))

// Check if there are more orders to load
const hasMoreOrders = computed(
  () => orders.value.meta.totalItems > displayedOrders.value.length
)

const router = useRouter()

// Navigate to order details
const goToOrder = (orderId: number) => {
  router.push(`/admin/orders/${orderId}`)
}

// Load more orders
const handleLoadMore = () => {
  const leftItemsCount = orders.value.meta.totalItems - limit.value
  if (leftItemsCount >= DEFAULT_LIMIT) {
    limit.value += DEFAULT_LIMIT
  } else if (leftItemsCount > 0) {
    limit.value += leftItemsCount
  }
}

// Format date using date-fns
const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'dd MMMM yyyy', { locale: ru })
}

// Get full name of the user
const getUserFullName = (user: { firstName: string; lastName: string }) => {
  return `${user.firstName} ${user.lastName}`
}

// Order status colors and formatted text
const ORDER_STATUS_COLOR: Record<string, string> = {
  completed: 'bg-green-100 text-green-800',
  pending: 'bg-yellow-100 text-yellow-800',
  cancelled: 'bg-red-100 text-red-800',
  // Add other statuses as needed
}

const ORDER_STATUS_FORMATTED: Record<string, string> = {
  completed: 'Завершен',
  pending: 'В ожидании',
  cancelled: 'Отменен',
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
