<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Запрос</TableHead>
				<TableHead>Склад</TableHead>
				<TableHead>Статус</TableHead>
				<TableHead>Дата создания</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="request in requests"
				:key="request.requestId"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="goToDetails(request.requestId)"
			>
				<TableCell class="font-medium">
					{{ request.requestId }}
				</TableCell>
				<TableCell>
					{{ request.warehouseName }}
				</TableCell>
				<TableCell>
					<p
						:class="[
              'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
              STATUS_COLOR[request.status]
            ]"
					>
						{{ STATUS_FORMATTED[request.status] }}
					</p>
				</TableCell>
				<TableCell>
					{{ new Date(request.createdAt).toLocaleDateString() }}
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
import { StoreStockRequestStatus, type StoreStockRequestResponse } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { useRouter } from 'vue-router'

const { requests } = defineProps<{ requests: StoreStockRequestResponse[] }>();

const router = useRouter();

const goToDetails = (requestId: number) => {
  router.push(`/admin/store-stock-requests/${requestId}`);
};

const STATUS_COLOR: Record<StoreStockRequestStatus, string> = {
  [StoreStockRequestStatus.CREATED]: 'bg-blue-100 text-blue-800',
  [StoreStockRequestStatus.IN_DELIVERY]: 'bg-yellow-100 text-yellow-800',
  [StoreStockRequestStatus.PROCESSED]: 'bg-gray-100 text-gray-800',
  [StoreStockRequestStatus.COMPLETED]: 'bg-green-100 text-green-800',
  [StoreStockRequestStatus.REJECTED]: 'bg-red-100 text-red-800',
};

const STATUS_FORMATTED: Record<StoreStockRequestStatus, string> = {
  [StoreStockRequestStatus.CREATED]: 'Создан',
  [StoreStockRequestStatus.IN_DELIVERY]: 'В доставке',
  [StoreStockRequestStatus.PROCESSED]: 'Обработан',
[  StoreStockRequestStatus.COMPLETED]: 'Завершён',
  [StoreStockRequestStatus.REJECTED]: 'Отклонён',
};
</script>
