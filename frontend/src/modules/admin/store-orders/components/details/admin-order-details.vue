<template>
	<div class="mx-auto max-w-7xl">
		<!-- Top Header with navigation -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="handleBack"
			>
				<ChevronLeft class="size-4" />
				<span class="sr-only">Назад</span>
			</Button>
			<div class="flex items-center gap-4">
				<h1 class="font-medium text-lg md:text-xl">
					Заказ #{{ orderDetails.displayNumber }} - {{ orderDetails.customerName }}
				</h1>
				<p
					:class="[
						'inline-flex w-fit items-center rounded-lg px-4 py-1 text-sm',
						ORDER_STATUS_COLOR[orderDetails.status]
					]"
				>
					{{ ORDER_STATUS_FORMATTED[orderDetails.status] }}
				</p>
			</div>
		</div>

		<!-- Main Grid Layout -->
		<div class="flex flex-col gap-4 md:grid md:grid-cols-6 mt-6">
			<!-- Left Column: Products and Payment -->
			<div class="flex flex-col gap-4 md:col-span-4 lg:col-span-3 xl:col-span-4">
				<!-- Suborders (Products) Card -->
				<Card>
					<CardHeader>
						<CardTitle>Подзаказы ({{ orderDetails.suborders.length }})</CardTitle>
						<CardDescription>
							Список подзаказов с данными о товаре, статусе и модификаторами
						</CardDescription>
					</CardHeader>
					<CardContent>
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Продукт</TableHead>
									<TableHead>Модификаторы</TableHead>
									<TableHead class="hidden md:table-cell p-4">Статус</TableHead>
									<TableHead class="text-right">Цена</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								<TableRow
									v-for="suborder in orderDetails.suborders"
									:key="suborder.id"
								>
									<TableCell class="flex items-center gap-4 font-medium">
										<LazyImage
											:src="suborder.storeProductSize.product.imageUrl"
											alt="Изображение товара"
											class="rounded-md w-12 h-12 object-cover"
										/>

										<RouterLink
											:to="{name: getRouteName('ADMIN_PRODUCT_DETAILS'), params: {id: suborder.storeProductSize.product.id}}"
											target="_blank"
											class="hover:text-primary underline underline-offset-4 transition-colors duration-300"
										>
											{{ suborder.storeProductSize.product.name }}
										</RouterLink>
									</TableCell>

									<TableCell>
										<span
											v-if="suborder.storeAdditives.length > 0"
											class="max-w-56 line-clamp-3"
										>
											{{ suborder.storeAdditives.map(a => a.name).join(', ') }}
										</span>
										<span v-else>Отсутсвуют</span>
									</TableCell>

									<TableCell class="hidden md:table-cell p-4">
										<p
											:class="[
												'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
												SUBORDER_STATUS_COLOR[suborder.status]
											]"
										>
											{{ SUBORDER_STATUS_FORMATTED[suborder.status] }}
										</p>
									</TableCell>

									<TableCell class="font-medium text-right">
										{{ formatPrice(suborder.price) }}
									</TableCell>
								</TableRow>
							</TableBody>
						</Table>
					</CardContent>
				</Card>

				<!-- Payment Card -->
				<Card>
					<CardHeader>
						<CardTitle>Оплата</CardTitle>
						<CardDescription>
							Платёжные данные: сумма подзаказов, скидки и итоговая стоимость
						</CardDescription>
					</CardHeader>
					<CardContent class="gap-4 grid">
						<div class="flex items-center">
							<p class="text-gray-500 text-base">
								Сумма подзаказов ({{ orderDetails.suborders.length }})
							</p>
							<div class="ml-auto font-medium">{{ formatPrice(subtotal) }}</div>
						</div>
						<div class="flex items-center">
							<p class="text-gray-500 text-base">Скидка</p>
							<div class="ml-auto font-medium">{{ formatPrice(discount) }}</div>
						</div>
						<Separator />
						<div class="flex items-center font-medium">
							<div>Итого</div>
							<div class="ml-auto font-medium">{{ formatPrice(orderDetails.total) }}</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Right Column: Customer & Contact Information -->
			<div class="flex flex-col gap-6 md:col-span-2 lg:col-span-3 xl:col-span-2">
				<Card>
					<!-- Покупатель -->
					<CardHeader>
						<div class="flex items-center gap-2">
							<User class="size-4" />
							<CardTitle class="font-medium">Покупатель</CardTitle>
						</div>
					</CardHeader>
					<CardContent>
						<p class="text-gray-500">
							{{ orderDetails.customerName || 'Не указан' }}
						</p>
					</CardContent>
					<Separator />

					<!-- Дата создания -->
					<CardHeader>
						<div class="flex items-center gap-2">
							<Calendar class="size-4" />
							<CardTitle class="font-medium">Дата создания</CardTitle>
						</div>
					</CardHeader>
					<CardContent>
						<p class="text-gray-500">
							{{ formatDate(orderDetails.createdAt) }}
						</p>
					</CardContent>
					<Separator />

					<!-- Дата завершения -->
					<CardHeader>
						<div class="flex items-center gap-2">
							<CalendarCheck class="size-4" />
							<CardTitle class="font-medium">Дата завершения</CardTitle>
						</div>
					</CardHeader>
					<CardContent>
						<p class="text-gray-500">
							{{ orderDetails.completedAt ? formatDate(orderDetails.completedAt) : 'Не завершен' }}
						</p>
					</CardContent>
					<Separator />

					<!-- Контактная информация -->
					<CardHeader>
						<div class="flex items-center gap-2">
							<Phone class="size-4" />
							<CardTitle class="font-medium">Контактная информация</CardTitle>
						</div>
					</CardHeader>
					<CardContent>
						<p class="text-gray-500">Телефон: отсутсвует</p>
					</CardContent>
					<Separator />

					<!-- Адрес доставки -->
					<CardHeader>
						<div class="flex items-center gap-2">
							<MapPin class="size-4" />
							<CardTitle class="font-medium">Адрес доставки</CardTitle>
						</div>
					</CardHeader>
					<CardContent>
						<div v-if="orderDetails.deliveryAddress">
							<p>{{ orderDetails.deliveryAddress.address }}</p>
						</div>
						<div v-else>
							<p class="text-gray-500">Адрес не указан</p>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { computed, toRefs } from 'vue'
import { useRouter } from 'vue-router'

// Import Shadcn-inspired UI components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { Separator } from '@/core/components/ui/separator'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'

// Import lucide icons
import { Calendar, CalendarCheck, ChevronLeft, MapPin, Phone, User } from 'lucide-vue-next'

// Import types and order models
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { OrderStatus, SubOrderStatus, type OrderDetailsDTO } from '@/modules/admin/store-orders/models/orders.models'

// ------ Props & Data Setup ------
const props = defineProps<{ orderDetails: OrderDetailsDTO }>()
const {orderDetails} = toRefs(props)
const router = useRouter()

// ------ Order & Suborder Status Mappings ------

const ORDER_STATUS_COLOR: Record<OrderStatus, string> = {
	[OrderStatus.WAITING_FOR_PAYMENT]: 'bg-pink-100 text-pink-800 border border-pink-200',
	[OrderStatus.PENDING]: 'bg-purple-100 text-purple-800 border border-purple-200',
	[OrderStatus.PREPARING]: 'bg-yellow-100 text-yellow-800 border border-yellow-200',
	[OrderStatus.COMPLETED]: 'bg-green-200 text-green-900 border border-green-200',
	[OrderStatus.IN_DELIVERY]: 'bg-orange-100 text-orange-800 border border-orange-200',
	[OrderStatus.DELIVERED]: 'bg-green-100 text-green-800 border border-green-200',
	[OrderStatus.CANCELLED]: 'bg-red-100 text-red-800 border border-red-200'
}

const ORDER_STATUS_FORMATTED: Record<OrderStatus, string> = {
	[OrderStatus.WAITING_FOR_PAYMENT]: 'Ожидание оплаты',
	[OrderStatus.PENDING]: 'В ожидании',
	[OrderStatus.PREPARING]: 'Готовится',
	[OrderStatus.COMPLETED]: 'Завершен',
	[OrderStatus.IN_DELIVERY]: 'Доставляется',
	[OrderStatus.DELIVERED]: 'Доставлен',
	[OrderStatus.CANCELLED]: 'Отменен'
}

const SUBORDER_STATUS_COLOR: Record<SubOrderStatus, string> = {
	[SubOrderStatus.PENDING]: 'bg-purple-100 text-purple-800',
	[SubOrderStatus.PREPARING]: 'bg-yellow-100 text-yellow-800',
	[SubOrderStatus.COMPLETED]: 'bg-green-200 text-green-900'
}

const SUBORDER_STATUS_FORMATTED: Record<SubOrderStatus, string> = {
	[SubOrderStatus.PENDING]: 'В ожидании',
	[SubOrderStatus.PREPARING]: 'Готовится',
	[SubOrderStatus.COMPLETED]: 'Завершен'
}

// ------ Computed Properties ------

// Compute the subtotal by summing prices of all suborders
const subtotal = computed(() =>
	orderDetails.value.suborders.reduce((sum, suborder) => sum + suborder.price, 0)
)
// For this example, discount is assumed to be zero
const discount = computed(() => 0)

// Format creation/complete dates
const formatDate = (dateString: Date | string) => {
  return format(new Date(dateString), 'dd MMMM yyyy, HH:mm', { locale: ru })
}

// ------ Navigation Handlers ------
function handleBack() {
	router.back()
}
</script>

<style scoped>
/* Additional scoped styling if needed */
</style>
