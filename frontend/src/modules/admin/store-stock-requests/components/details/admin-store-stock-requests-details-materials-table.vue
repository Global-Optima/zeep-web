<template>
	<Card>
		<CardHeader>
			<div class="flex justify-between items-start gap-4">
				<div>
					<CardTitle>Список материалов</CardTitle>
					<CardDescription class="mt-2">
						Ниже представлена таблица с материалами, переданная в компонент через props.
					</CardDescription>
				</div>

				<Button
					v-if="isEditable"
					size="icon"
					variant="ghost"
					@click="onUpdateClick"
				>
					<Pencil class="w-5 h-5 text-gray-600" />
				</Button>
			</div>
		</CardHeader>

		<CardContent>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Материал</TableHead>
						<TableHead>Количество</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow
						v-for="(item, index) in stockRequest.items"
						:key="index"
					>
						<TableCell>{{ item.name }}</TableCell>
						<TableCell>{{ item.quantity }}</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/core/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'
import { type StoreStockRequestResponse, StoreStockRequestStatus } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { Pencil } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const showEditButtonStatuses: StoreStockRequestStatus[] = [StoreStockRequestStatus.CREATED]

const {stockRequest} = defineProps<{
  stockRequest: StoreStockRequestResponse
}>()

const isEditable = computed(() => showEditButtonStatuses.includes(stockRequest.status))

const router = useRouter()

const onUpdateClick = () => {
  router.push(`/admin/store-stock-requests/${stockRequest.requestId}/update`)
}
</script>

<style scoped></style>
