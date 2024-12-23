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
				<!-- Request ID -->
				<div>
					<p class="text-muted-foreground text-sm">Номер заявки</p>
					<p>{{ request.requestId }}</p>
				</div>

				<!-- Store Name -->
				<div>
					<p class="text-muted-foreground text-sm">Магазин</p>
					<p>{{ request.storeName }}</p>
				</div>

				<!-- Warehouse Name -->
				<div>
					<p class="text-muted-foreground text-sm">Склад</p>
					<p>{{ request.warehouseName }}</p>
				</div>

				<!-- Status -->
				<div>
					<p class="text-muted-foreground text-sm">Статус</p>
					<p>{{ statusFormatted }}</p>
				</div>

				<!-- Created At -->
				<div>
					<p class="text-muted-foreground text-sm">Дата создания</p>
					<p>{{ new Date(request.createdAt).toLocaleDateString() }}</p>
				</div>

				<!-- Updated At -->
				<div>
					<p class="text-muted-foreground text-sm">Дата обновления</p>
					<p>{{ new Date(request.updatedAt).toLocaleDateString() }}</p>
				</div>
			</div>
		</CardContent>

		<CardFooter>
			<div class="space-y-2 w-full">
				<Button
					v-if="props.request.status === StoreStockRequestStatus.PROCESSED"
					variant="destructive"
					:disabled="props.request.status !== 'PROCESSED'"
					@click="rejectRequest"
					class="w-full"
				>
					{{ props.request.status === 'PROCESSED' ? 'Отклонить заявку' : 'Отклонение недоступно' }}
				</Button>
				<Button
					:disabled="!isActionAllowed"
					@click="handleStatusChange"
					class="w-full"
				>
					{{ props.request.status === 'PROCESSED' ? 'Отправить в доставку' : 'Изменение недоступно' }}
				</Button>
			</div>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  StoreStockRequestStatus,
  type StoreStockRequestResponse,
} from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { computed } from 'vue'

const props = defineProps<{ request: StoreStockRequestResponse }>();
const emit = defineEmits<{ (e: 'update:status', newStatus: StoreStockRequestStatus): void }>();

// Helper to format status
const statusLabels: Record<StoreStockRequestStatus, string> = {
  CREATED: 'Создана',
  PROCESSED: 'Запрос отправлен',
  IN_DELIVERY: 'В доставке',
  COMPLETED: 'Завершена',
  REJECTED: 'Отклонена',
};

const statusFormatted = computed(() => statusLabels[props.request.status]);


// Determine if action is allowed
const isActionAllowed = computed(() => {
  return props.request.status === StoreStockRequestStatus.PROCESSED;
});

// Handle status change to IN_DELIVERY
function handleStatusChange() {
  if (props.request.status === StoreStockRequestStatus.PROCESSED) {
    emit('update:status', StoreStockRequestStatus.IN_DELIVERY);
  }
}

// Handle rejection of the request
function rejectRequest() {
  if (props.request.status === StoreStockRequestStatus.PROCESSED) {
    emit('update:status', StoreStockRequestStatus.REJECTED);
  }
}
</script>

<style scoped>
.text-muted-foreground {
  color: #6b7280; /* Tailwind muted text */
}
</style>
