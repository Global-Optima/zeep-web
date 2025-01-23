<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Доставка №{{ delivery.id }} - {{ delivery.supplier.name }}
			</h1>
		</div>

		<!-- Delivery Info -->
		<Card>
			<CardHeader>
				<CardTitle>Информация о доставке</CardTitle>
				<CardDescription>Основные данные о доставке</CardDescription>
			</CardHeader>
			<CardContent>
				<ul class="space-y-2">
					<li>
						<span>Склад: </span> <span class="font-medium">{{ delivery.warehouse.name}}</span>
					</li>
					<li>
						<span>Количество: </span>
						<span class="font-medium">{{ delivery.materials.length }}</span>
					</li>
					<li>
						<span>Поставщик: </span>
						<span class="font-medium">{{ delivery.supplier.name }}</span>
					</li>
					<li>
						<span>Склад: </span> <span class="font-medium">{{ delivery.warehouse.name }}</span>
					</li>
					<li>
						<span>Дата доставки: yyyy</span>
						<span class="font-medium">{{ formatDate(delivery.deliveryDate) }}</span>
					</li>
				</ul>
			</CardContent>
		</Card>

		<!-- Stock Materials Table -->
		<Card>
			<CardHeader>
				<CardTitle>Материалы в доставке</CardTitle>
				<CardDescription>Список материалов, включенных в доставку</CardDescription>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название материала</TableHead>
							<TableHead>Упаковка</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead>Срок годности</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="material in delivery.materials"
							:key="material.stockMaterial.id"
						>
							<TableCell class="py-4 font-medium">{{ material.stockMaterial.name }}</TableCell>
							<TableCell>{{ material.stockMaterial.unit.name }}</TableCell>
							<TableCell>{{ material.quantity }}</TableCell>
							<TableCell>{{ format(material.expirationDate, "dd.MM.yyyy") }}</TableCell>
						</TableRow>
						<TableRow v-if="delivery.materials.length === 0">
							<TableCell
								colspan="3"
								class="py-4 text-center text-muted-foreground"
							>
								Материалы отсутствуют
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import type { WarehouseDeliveryDTO } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { format } from 'date-fns'
import { ChevronLeft } from 'lucide-vue-next'

// Props

const { delivery } = defineProps<{
	delivery: WarehouseDeliveryDTO
}>()

const emit = defineEmits<{
	(e: 'onCancel'): void
}>()

// Utility Function: Format Date
function formatDate(date: Date): string {
	return new Date(date).toLocaleDateString('ru-RU', {
		day: '2-digit',
		month: '2-digit',
		year: 'numeric',
	})
}

const onCancel = () => {
	emit('onCancel')
}
</script>
