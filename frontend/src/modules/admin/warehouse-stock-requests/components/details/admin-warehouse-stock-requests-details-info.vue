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
				<!-- Request details remain the same -->
				<div>
					<p class="text-muted-foreground text-sm">Номер заявки</p>
					<p>{{ request.requestId }}</p>
				</div>
				<div>
					<p class="text-muted-foreground text-sm">Магазин</p>
					<p>{{ request.store.name }}</p>
				</div>
				<div>
					<p class="text-muted-foreground text-sm">Склад</p>
					<p>{{ request.warehouse.name }}</p>
				</div>
				<div>
					<p class="text-muted-foreground text-sm">Статус</p>
					<p>{{ statusFormatted }}</p>
				</div>
				<div>
					<p class="text-muted-foreground text-sm">Дата создания</p>
					<p>{{ new Date(request.createdAt).toLocaleDateString() }}</p>
				</div>
				<div>
					<p class="text-muted-foreground text-sm">Дата обновления</p>
					<p>{{ new Date(request.updatedAt).toLocaleDateString() }}</p>
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

const props = defineProps<{ request: StockRequestResponse }>()

const {currentEmployee} = useEmployeeAuthStore()
const userRole = computed(() => currentEmployee?.role)

const statusFormatted = computed(() => {
  return STOCK_REQUEST_STATUS_FORMATTED[props.request.status]
})
</script>

<style scoped></style>
