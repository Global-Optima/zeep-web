<template>
	<Card>
		<CardHeader>
			<CardTitle>Список материалов</CardTitle>
			<CardDescription class="mt-2"> Материалы представленные в заказе </CardDescription>
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
						<TableCell>{{ item.quantity }}</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
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
import { type StockRequestResponse } from '@/modules/admin/store-stock-requests/models/stock-requests.model'


const {stockRequest} = defineProps<{
  stockRequest: StockRequestResponse
}>()
</script>

<style scoped></style>
