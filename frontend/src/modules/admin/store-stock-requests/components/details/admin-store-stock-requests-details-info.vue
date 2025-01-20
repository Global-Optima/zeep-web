<template>
	<Card>
		<CardHeader>
			<div>
				<CardTitle>Информация о заявке</CardTitle>
				<CardDescription class="mt-2">
					Подробная информация о заявке на материалы.
				</CardDescription>
			</div>
		</CardHeader>

		<CardContent>
			<div class="space-y-4">
				<!-- Iterate over predefined requestDetails array -->
				<div
					v-for="detail in requestDetails"
					:key="detail.label"
				>
					<p class="mb-1 text-muted-foreground text-sm">{{ detail.label }}</p>
					<p v-if="detail.value">
						{{ detail.value }}
					</p>
					<p v-else>Отсутствует</p>
				</div>
			</div>
		</CardContent>

		<CardFooter>
			<div
				v-if="userRole"
				class="w-full"
			>
				<AdminStockRequestsActions
					:status="request.status"
					:request="request"
					:role="userRole"
				/>
			</div>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import AdminStockRequestsActions from '@/modules/admin/store-stock-requests/components/details/admin-stock-requests-actions.vue'
import { type StockRequestResponse, STOCK_REQUEST_STATUS_FORMATTED } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { computed } from 'vue'

// Props
const props = defineProps<{ request: StockRequestResponse }>()

// Current user role
const { currentEmployee } = useEmployeeAuthStore()
const userRole = computed(() => currentEmployee?.role)

// Predefined array for request details
const requestDetails = computed(() => [
  { label: 'Номер заявки', value: props.request.requestId },
  { label: 'Магазин', value: props.request.store.name },
  { label: 'Склад', value: props.request.warehouse.name },
  { label: 'Статус', value: STOCK_REQUEST_STATUS_FORMATTED[props.request.status] },
  {
    label: 'Дата создания',
    value: new Date(props.request.createdAt).toLocaleDateString('ru-RU'),
  },
  {
    label: 'Дата обновления',
    value: new Date(props.request.updatedAt).toLocaleDateString('ru-RU'),
  },
  { label: 'Комментарий от заказчика', value: props.request.storeComment },
  { label: 'Комментарий от Склада', value: props.request.warehouseComment },
])
</script>

<style scoped></style>
