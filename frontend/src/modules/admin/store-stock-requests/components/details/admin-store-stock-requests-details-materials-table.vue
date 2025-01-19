<template>
	<Card>
		<CardHeader>
			<div class="flex justify-between items-start gap-4">
				<div>
					<CardTitle>Список материалов</CardTitle>
					<CardDescription class="mt-2"> Материалы представленные в заказе </CardDescription>
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
						v-for="(item, index) in stockRequest.stockMaterials"
						:key="index"
					>
						<TableCell>{{ item.stockMaterial.name }}</TableCell>
						<TableCell>{{ item.packageMeasures.quantity }}</TableCell>
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
import { type StockRequestResponse, StockRequestStatus } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { Pencil } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const showEditButtonStatuses: StockRequestStatus[] = [StockRequestStatus.CREATED, StockRequestStatus.REJECTED_BY_WAREHOUSE]

const {stockRequest} = defineProps<{
  stockRequest: StockRequestResponse
}>()

const isEditable = computed(() => showEditButtonStatuses.includes(stockRequest.status))

const router = useRouter()

const onUpdateClick = () => {
  router.push(`/admin/store-stock-requests/${stockRequest.requestId}/update`)
}
</script>

<style scoped></style>
