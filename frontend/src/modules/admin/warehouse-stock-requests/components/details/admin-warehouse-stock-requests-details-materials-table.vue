<template>
	<Card>
		<CardHeader>
			<div>
				<CardTitle>Список материалов</CardTitle>
				<CardDescription class="mt-2">
					Ниже представлена таблица с материалами, переданная в компонент через props.
				</CardDescription>
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
						<TableCell>{{ item.name }}</TableCell>
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
import type { StoreStockRequestItemResponse } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'

defineProps<{
  items: StoreStockRequestItemResponse[]
}>()
</script>

<style scoped></style>
