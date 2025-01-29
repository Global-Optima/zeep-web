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
					{{ request.warehouse.name }}
				</TableCell>
				<TableCell>
					<p
						:class="[
              'inline-flex w-fit items-center rounded-md px-2.5 py-1 text-xs',
              STOCK_REQUEST_STATUS_COLOR[request.status]
            ]"
					>
						{{ STOCK_REQUEST_STATUS_FORMATTED[request.status] }}
					</p>
				</TableCell>
				<TableCell>
					{{format(request.createdAt, "dd.MM.yyyy hh:mm", {locale: ru})}}
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
import { STOCK_REQUEST_STATUS_COLOR, STOCK_REQUEST_STATUS_FORMATTED, type StockRequestResponse } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { useRouter } from 'vue-router'

const { requests } = defineProps<{ requests: StockRequestResponse[] }>();

const router = useRouter();

const goToDetails = (requestId: number) => {
  router.push(`/admin/store-stock-requests/${requestId}`);
};
</script>
