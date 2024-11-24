<template>
	<div class="relative bg-gray-100 pt-safe w-full min-h-screen">
		<!-- Header: Order Status Selector -->
		<header
			class="top-0 left-0 fixed flex justify-center items-center bg-white p-4 border-b w-full overflow-x-auto no-scrollbar"
		>
			<button
				v-for="status in statuses"
				:key="status.label"
				@click="onSelectStatus(status)"
				:class="[
          'flex items-center gap-2 px-5 py-2 rounded-xl text-lg whitespace-nowrap',
          status.label === selectedStatus.label ? 'bg-primary text-primary-foreground' : ''
        ]"
			>
				<p>{{ status.label }}</p>
				<p
					:class="[
            'bg-gray-100 px-2 py-1 rounded-sm text-black text-xs',
            status.label === selectedStatus.label ? 'bg-green-700 text-primary-foreground' : ''
          ]"
				>
					{{ status.count }}
				</p>
			</button>
		</header>

		<!-- Main Layout -->
		<div class="grid grid-cols-4 bg-gray-100 mt-[77px] w-full min-h-full">
			<!-- Section 1: Orders -->
			<section class="col-span-1 p-2 border-r h-full">
				<p class="py-1 font-medium text-center">Заказы</p>

				<div
					v-if="filteredOrders.length > 0"
					class="flex flex-col gap-2 mt-2 overflow-y-auto no-scrollbar"
					ref="ordersSection"
				>
					<div
						v-for="order in filteredOrders"
						:key="order.id"
						@click="selectOrder(order)"
						:class="[
              'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
              selectedOrder?.id === order.id ? 'border-primary' : 'border-transparent'
            ]"
					>
						<div>
							<div class="flex items-center gap-2 text-lg">
								<p class="font-medium">{{ order.customerName }}</p>
								<p class="text-gray-500">#{{ order.id }}</p>
							</div>
							<div class="mt-1 text-gray-700 text-sm">
								<span>{{ order.type === 'Delivery' ? 'Доставка' : 'Кафе' }}</span
								>,
								<span>{{ order.suborders.length }} шт.</span>
							</div>
						</div>
						<div>
							<template v-if="order.status === 'Active'">
								<p class="text-blue-600">{{ order.eta }}</p>
							</template>
							<template v-else-if="order.status === 'In Delivery'">
								<Truck class="w-5 h-5 text-yellow-500" />
							</template>
							<template v-else-if="order.status === 'Completed'">
								<Check class="w-5 h-5 text-green-500" />
							</template>
						</div>
					</div>
				</div>

				<p
					v-else
					class="mt-2 text-center text-gray-400 text-sm"
				>
					Список заказов пуст
				</p>
			</section>

			<!-- Section 2: Suborders -->
			<section class="col-span-1 p-2 border-r h-full">
				<p class="py-1 font-medium text-center">Подзаказы</p>

				<div
					v-if="selectedOrder"
					class="flex flex-col flex-1 gap-2 mt-2 overflow-y-auto no-scrollbar"
					ref="subordersSection"
				>
					<div
						v-for="suborder in selectedOrder.suborders"
						:key="suborder.id"
						@click="selectSuborder(suborder)"
						:class="[
              'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
              selectedSuborder?.id === suborder.id ? 'border-primary' : 'border-transparent'
            ]"
					>
						<div>
							<p class="font-medium text-lg">{{ suborder.productName }}</p>
							<p class="line-clamp-2 text-gray-700 text-sm">
								{{ suborder.toppings.join(', ') }}
							</p>
						</div>
						<div>
							<p
								v-if="suborder.status === 'In Progress'"
								class="text-blue-600"
							>
								{{ suborder.prepTime }}
							</p>
							<Check
								v-else
								class="w-5 h-5 text-green-500"
							/>
						</div>
					</div>
				</div>

				<p
					v-else
					class="mt-2 text-center text-gray-400 text-sm"
				>
					Выберите заказ
				</p>
			</section>

			<!-- Section 3: Suborder Details -->
			<section class="col-span-2 p-2 h-full">
				<p class="py-1 font-medium text-center">Детали подзаказа</p>
				<div
					v-if="selectedSuborder"
					class="flex flex-col gap-4 bg-white mt-2 p-4 rounded-xl overflow-y-auto"
					ref="detailsSection"
				>
					<div>
						<p class="font-medium text-xl">{{ selectedSuborder.productName }}</p>
						<!-- Toppings List -->
						<ul
							v-if="selectedSuborder.toppings.length > 0"
							class="space-y-1 mt-2"
						>
							<li
								v-for="(topping, index) in selectedSuborder.toppings"
								:key="index"
								class="flex items-center"
							>
								<Plus class="mr-2 w-4 h-4 text-gray-500" />
								<span class="text-gray-700">{{ topping }}</span>
							</li>
						</ul>
						<p
							v-else
							class="mt-2 text-gray-700"
						>
							Без топпингов
						</p>
					</div>
					<div>
						<p class="font-medium text-lg">Комментарий</p>
						<p class="mt-1 text-gray-700">
							{{ selectedSuborder.comments || 'Стандартное приготовление' }}
						</p>
					</div>
					<div>
						<p class="font-medium text-lg">Время приготовления</p>
						<p class="mt-1 text-gray-700">
							{{ selectedSuborder.prepTime || 'Не указано' }}
						</p>
					</div>
					<div class="flex items-center gap-2 mt-4">
						<button
							class="p-4 border rounded-xl"
							:disabled="selectedSuborder.status === 'Done'"
							@click="printQrCode"
						>
							<Printer class="w-6 h-6" />
						</button>
						<button
							@click="toggleSuborderStatus(selectedSuborder)"
							:disabled="selectedSuborder.status === 'Done'"
							:class="[
                'flex-1 px-4 py-4 rounded-xl text-primary-foreground',
                selectedSuborder.status === 'Done'
                  ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                  : 'bg-primary'
              ]"
						>
							{{ selectedSuborder.status === 'Done' ? 'Выполнено' : 'Выполнить' }}
						</button>
					</div>
				</div>

				<p
					v-else
					class="mt-2 text-center text-gray-400 text-sm"
				>
					Выберите подзаказ
				</p>
			</section>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Check, Plus, Printer, Truck } from 'lucide-vue-next'
import { computed, nextTick, ref, watch } from 'vue'

/**
 * Interfaces for typing
 */
interface Suborder {
  id: number
  productName: string
  toppings: string[]
  status: 'In Progress' | 'Done'
  comments?: string
  prepTime: string
}

interface Order {
  id: number
  customerName: string
  customerEmail: string
  details: string
  eta: string
  suborders: Suborder[]
  status: 'Active' | 'Completed' | 'In Delivery'
  type: 'Delivery' | 'In-Store'
}

interface Status {
  label: string
  count: number
}

const scrollToTop = async () => {
  await nextTick();
  if (window && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};


/**
 * Mock data with different statuses and order types
 */
 const orders = ref<Order[]>([
  // Active Orders
  {
    id: 3302,
    customerName: 'Алексей Иванов',
    customerEmail: 'alexey.ivanov@example.com',
    details: 'Доставка, 3 шт',
    eta: '~5 мин',
    status: 'Active',
    type: 'Delivery',
    suborders: [
      {
        id: 1,
        productName: 'Капучино',
        toppings: ['Молоко', 'Сахар', 'Корица'],
        status: 'In Progress',
        comments: 'Добавить больше молока',
        prepTime: '3 мин',
      },
      {
        id: 2,
        productName: 'Латте',
        toppings: ['Ванильный сироп'],
        status: 'In Progress',
        comments: 'Без пены',
        prepTime: '4 мин',
      },
      {
        id: 3,
        productName: 'Эспрессо',
        toppings: [],
        status: 'In Progress',
        prepTime: '2 мин',
      },
    ],
  },
  {
    id: 3303,
    customerName: 'Мария Петрова',
    customerEmail: 'maria.petrova@example.com',
    details: 'Кафе, 2 шт',
    eta: '~2 мин',
    status: 'Active',
    type: 'In-Store',
    suborders: [
      {
        id: 4,
        productName: 'Американо',
        toppings: ['Корица'],
        status: 'In Progress',
        prepTime: '2 мин',
      },
      {
        id: 5,
        productName: 'Мокачино',
        toppings: ['Шоколадный сироп'],
        status: 'In Progress',
        comments: 'Двойная порция сиропа',
        prepTime: '5 мин',
      },
    ],
  },
  {
    id: 3306,
    customerName: 'Иван Кузнецов',
    customerEmail: 'ivan.kuznetsov@example.com',
    details: 'Доставка, 2 шт',
    eta: '~10 мин',
    status: 'Active',
    type: 'Delivery',
    suborders: [
      {
        id: 6,
        productName: 'Эспрессо',
        toppings: [],
        status: 'In Progress',
        prepTime: '2 мин',
      },
      {
        id: 7,
        productName: 'Американо',
        toppings: ['Молоко'],
        status: 'In Progress',
        prepTime: '3 мин',
      },
    ],
  },
  {
    id: 3307,
    customerName: 'Ольга Морозова',
    customerEmail: 'olga.morozova@example.com',
    details: 'Кафе, 1 шт',
    eta: '~1 мин',
    status: 'Active',
    type: 'In-Store',
    suborders: [
      {
        id: 8,
        productName: 'Чай чёрный',
        toppings: [],
        status: 'In Progress',
        prepTime: '1 мин',
      },
    ],
  },
  {
    id: 3310,
    customerName: 'Павел Фёдоров',
    customerEmail: 'pavel.fedorov@example.com',
    details: 'Доставка, 1 шт',
    eta: '~7 мин',
    status: 'Active',
    type: 'Delivery',
    suborders: [
      {
        id: 9,
        productName: 'Фраппучино',
        toppings: ['Ванильный сироп'],
        status: 'In Progress',
        prepTime: '6 мин',
      },
    ],
  },

  // In Delivery Orders
  {
    id: 3305,
    customerName: 'Елена Соколова',
    customerEmail: 'elena.sokolova@example.com',
    details: 'Кафе, 2 шт',
    eta: 'В пути',
    status: 'In Delivery',
    type: 'In-Store',
    suborders: [
      {
        id: 10,
        productName: 'Чай латте',
        toppings: ['Мёд'],
        status: 'Done',
        prepTime: '5 мин',
      },
      {
        id: 11,
        productName: 'Матча латте',
        toppings: [],
        status: 'Done',
        prepTime: '5 мин',
      },
    ],
  },
  {
    id: 3308,
    customerName: 'Дмитрий Волков',
    customerEmail: 'dmitry.volkov@example.com',
    details: 'Доставка, 1 шт',
    eta: 'В пути',
    status: 'In Delivery',
    type: 'Delivery',
    suborders: [
      {
        id: 12,
        productName: 'Латте',
        toppings: ['Корица'],
        status: 'Done',
        prepTime: '4 мин',
      },
    ],
  },

  // Completed Orders
  {
    id: 3304,
    customerName: 'Сергей Смирнов',
    customerEmail: 'sergey.smirnov@example.com',
    details: 'Доставка, 1 шт',
    eta: 'Доставлен',
    status: 'Completed',
    type: 'Delivery',
    suborders: [
      {
        id: 13,
        productName: 'Фраппучино',
        toppings: ['Карамельный сироп'],
        prepTime: '5 мин',
        status: 'Done',
      },
    ],
  },
  {
    id: 3309,
    customerName: 'Наталья Васильева',
    customerEmail: 'natalia.vasileva@example.com',
    details: 'Кафе, 2 шт',
    eta: 'Доставлен',
    status: 'Completed',
    type: 'In-Store',
    suborders: [
      {
        id: 14,
        productName: 'Мокачино',
        toppings: ['Шоколадный сироп'],
        status: 'Done',
        prepTime: '5 мин',
      },
      {
        id: 15,
        productName: 'Капучино',
        toppings: [],
        status: 'Done',
        prepTime: '3 мин',
      },
    ],
  },
  {
    id: 3311,
    customerName: 'Анна Попова',
    customerEmail: 'anna.popova@example.com',
    details: 'Кафе, 1 шт',
    eta: 'Доставлен',
    status: 'Completed',
    type: 'In-Store',
    suborders: [
      {
        id: 16,
        productName: 'Чай зелёный',
        toppings: [],
        status: 'Done',
        prepTime: '2 мин',
      },
    ],
  },
  {
    id: 3312,
    customerName: 'Максим Васильев',
    customerEmail: 'maxim.vasiliev@example.com',
    details: 'Доставка, 2 шт',
    eta: 'Доставлен',
    status: 'Completed',
    type: 'Delivery',
    suborders: [
      {
        id: 17,
        productName: 'Американо',
        toppings: [],
        status: 'Done',
        prepTime: '3 мин',
      },
      {
        id: 18,
        productName: 'Эспрессо',
        toppings: [],
        status: 'Done',
        prepTime: '2 мин',
      },
    ],
  },
]);


/**
 * Reactive states
 */
const statuses = ref<Status[]>([
  { label: 'Все', count: orders.value.length },
  { label: 'Активные', count: orders.value.filter(o => o.status === 'Active').length },
  { label: 'Завершенные', count: orders.value.filter(o => o.status === 'Completed').length },
  { label: 'В доставке', count: orders.value.filter(o => o.status === 'In Delivery').length }
])

const selectedStatus = ref<Status>(statuses.value[0])
const selectedOrder = ref<Order | null>(null)
const selectedSuborder = ref<Suborder | null>(null)

/**
 * Computed values
 */
const filteredOrders = computed(() => {
  if (selectedStatus.value.label === 'Все') {
    return orders.value
  }
  if (selectedStatus.value.label === 'Активные') {
    return orders.value.filter(order => order.status === 'Active')
  }
  if (selectedStatus.value.label === 'Завершенные') {
    return orders.value.filter(order => order.status === 'Completed')
  }
  if (selectedStatus.value.label === 'В доставке') {
    return orders.value.filter(order => order.status === 'In Delivery')
  }
  return []
})

/**
 * Methods
 */
const onSelectStatus = async (status: Status) => {
  selectedStatus.value = status
  selectedOrder.value = null
  selectedSuborder.value = null

  await scrollToTop()
}

const selectOrder = async (order: Order) => {
  if (selectedOrder.value?.id === order.id) return
  selectedOrder.value = order
  selectedSuborder.value = null

  await scrollToTop()

}

const selectSuborder = async (suborder: Suborder) => {
  if (selectedSuborder.value?.id === suborder.id) return
  selectedSuborder.value = suborder

  await scrollToTop()
}

/**
 * Toggle suborder status between 'In Progress' and 'Done'
 */
const toggleSuborderStatus = (suborder: Suborder) => {
  if (suborder.status === 'Done') return
  suborder.status = 'Done'

  // Check if all suborders are done to mark the order as completed
  const allDone = selectedOrder.value?.suborders.every(so => so.status === 'Done')
  if (allDone && selectedOrder.value?.status === 'Active') {
    selectedOrder.value.status = 'Completed'
    updateStatusCounts()
    // Unselect order and suborder when order is completed
    selectedOrder.value = null
    selectedSuborder.value = null
  }
}

const printQrCode = () => {
  console.log('Print QR Code')
}


/**
 * Watchers
 */
watch(
  () => orders.value,
  () => {
    updateStatusCounts()
  },
  { deep: true }
)

/**
 * Update counts in statuses based on orders
 */
const updateStatusCounts = () => {
  statuses.value = statuses.value.map(status => {
    if (status.label === 'Все') {
      return { ...status, count: orders.value.length }
    }
    if (status.label === 'Активные') {
      return { ...status, count: orders.value.filter(o => o.status === 'Active').length }
    }
    if (status.label === 'Завершенные') {
      return { ...status, count: orders.value.filter(o => o.status === 'Completed').length }
    }
    if (status.label === 'В доставке') {
      return { ...status, count: orders.value.filter(o => o.status === 'In Delivery').length }
    }
    return status
  })
}
</script>

<style scoped></style>
