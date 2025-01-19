<template>
	<Card>
		<CardHeader>
			<div>
				<CardTitle>Список материалов</CardTitle>
				<CardDescription class="mt-2"> Материалы представленные в заказе </CardDescription>
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
						v-for="(item, index) in items"
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
import type { StockRequestMaterial } from '@/modules/admin/store-stock-requests/models/stock-requests.model'

defineProps<{
  items: StockRequestMaterial[]
}>()
</script>

<style scoped></style>
