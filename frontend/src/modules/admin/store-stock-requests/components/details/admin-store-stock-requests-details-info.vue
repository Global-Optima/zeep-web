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
					<p>{{ request.store.name }}</p>
				</div>

				<!-- Warehouse Name -->
				<div>
					<p class="text-muted-foreground text-sm">Склад</p>
					<p>{{ request.warehouse.name }}</p>
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
			<Button
				:disabled="!isActionAllowed"
				@click="handleStatusChange"
				class="w-full"
			>
				{{ getButtonLabel }}
			</Button>
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
import { STOCK_REQUEST_STATUS_FORMATTED, type StockRequestResponse, StockRequestStatus } from '@/modules/admin/store-stock-requests/models/stock-requests.model'

import { computed } from 'vue'

const props = defineProps<{ request: StockRequestResponse }>();
const emit = defineEmits<{ (e: 'update:status', newStatus: StockRequestStatus): void }>();

const statusLabels: Record<StockRequestStatus, string> = STOCK_REQUEST_STATUS_FORMATTED

const statusFormatted = computed(() => statusLabels[props.request.status]);

const getButtonLabel = computed(() => {
  switch (props.request.status) {
    case 'CREATED':
      return 'Отправить на склад';
    case 'IN_DELIVERY':
      return 'Завершить заявку';
      case 'PROCESSED':
      return 'Заявка отправлена на склад';
    case 'COMPLETED':
      return 'Заявка завершена';
    default:
      return '';
  }
});

const isActionAllowed = computed(() => {
  return props.request.status === StockRequestStatus.CREATED || props.request.status === StockRequestStatus.IN_DELIVERY;
});

function handleStatusChange() {
  let newStatus: StockRequestStatus | null = null;

  if (props.request.status === StockRequestStatus.CREATED) {
    newStatus = StockRequestStatus.PROCESSED;
  } else if (props.request.status === StockRequestStatus.IN_DELIVERY) {
    newStatus = StockRequestStatus.COMPLETED;
  }

  if (newStatus) {
    emit('update:status', newStatus);
  }
}
</script>

<style scoped>
.text-muted-foreground {
  color: #6b7280; /* Tailwind muted text */
}
</style>
