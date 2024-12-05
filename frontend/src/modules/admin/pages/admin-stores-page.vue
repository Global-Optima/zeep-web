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
					class="hover:bg-gray-50 h-12 cursor-pointer"
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
      id: 4,
      createdAt: '2023-10-04T14:45:12Z',
      user: {
        firstName: 'Алексей',
        lastName: 'Смирнов',
      },
      totalPrice: 1520.50,
      isPaid: true,
      status: 'completed',
      notes: 'Позвонить за час до доставки',
    },
    {
      id: 5,
      createdAt: '2023-10-05T10:20:30Z',
      user: {
        firstName: 'Ольга',
        lastName: 'Кузнецова',
      },
      totalPrice: 250.0,
      isPaid: false,
      status: 'pending',
      notes: '',
    },
    {
      id: 6,
      createdAt: '2023-10-06T09:05:15Z',
      user: {
        firstName: 'Дмитрий',
        lastName: 'Попов',
      },
      totalPrice: 3100.75,
      isPaid: true,
      status: 'completed',
      notes: 'Оставить у двери',
    },
    {
      id: 7,
      createdAt: '2023-10-07T18:25:00Z',
      user: {
        firstName: 'Елена',
        lastName: 'Васильева',
      },
      totalPrice: 520.90,
      isPaid: false,
      status: 'pending',
      notes: 'Уточнить время доставки',
    },
    {
      id: 8,
      createdAt: '2023-10-08T13:45:45Z',
      user: {
        firstName: 'Михаил',
        lastName: 'Соколов',
      },
      totalPrice: 870.40,
      isPaid: true,
      status: 'completed',
      notes: '',
    },
    {
      id: 9,
      createdAt: '2023-10-09T16:00:00Z',
      user: {
        firstName: 'Анна',
        lastName: 'Зайцева',
      },
      totalPrice: 1445.20,
      isPaid: false,
      status: 'pending',
      notes: 'Привезти до 12:00',
    },
    {
      id: 10,
      createdAt: '2023-10-10T11:15:30Z',
      user: {
        firstName: 'Игорь',
        lastName: 'Лебедев',
      },
      totalPrice: 3400.0,
      isPaid: true,
      status: 'completed',
      notes: '',
    },
    {
      id: 11,
      createdAt: '2023-10-11T19:30:00Z',
      user: {
        firstName: 'Виктория',
        lastName: 'Морозова',
      },
      totalPrice: 645.30,
      isPaid: false,
      status: 'pending',
      notes: '',
    },
    {
      id: 12,
      createdAt: '2023-10-12T08:45:00Z',
      user: {
        firstName: 'Сергей',
        lastName: 'Киселёв',
      },
      totalPrice: 2450.75,
      isPaid: true,
      status: 'completed',
      notes: 'Позвонить по прибытии',
    },
    {
      id: 13,
      createdAt: '2023-10-13T15:10:20Z',
      user: {
        firstName: 'Татьяна',
        lastName: 'Фёдорова',
      },
      totalPrice: 305.60,
      isPaid: false,
      status: 'pending',
      notes: 'Без сдачи',
    },
  ],
  meta: {
    totalItems: 20, // Total number of orders (for pagination)
  },
});


const limit = ref(DEFAULT_LIMIT)

// Computed property for displayed orders based on limit
const displayedOrders = computed(() => orders.value.data.slice(0, limit.value))

const router = useRouter()

// Navigate to order details
const goToOrder = (orderId: number) => {
  router.push(`/admin/orders/${orderId}`)
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
