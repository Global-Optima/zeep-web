<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'

import type {
  SupplierDTO,
  SupplierMaterialResponse,
  UpdateSupplierMaterialDTO
} from '@/modules/admin/suppliers/models/suppliers.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

interface UpdateSupplierMaterialForm extends UpdateSupplierMaterialDTO {
  name: string
  unit: UnitDTO
  size: number
}

/** Props */
const { supplier, stockMaterials } = defineProps<{ supplier: SupplierDTO, stockMaterials: SupplierMaterialResponse[] }>()
const emit = defineEmits<{
  (e: 'on-cancel'): void
}>()


const mergedStockMaterialsTableData = ref<UpdateSupplierMaterialForm[]>(stockMaterials.map(m => ({
    stockMaterialId: m.stockMaterial.id,
    name: m.stockMaterial.name,
    unit: m.stockMaterial.unit,
    basePrice: m.basePrice,
    size: m.stockMaterial.size
  })))

/** Cancel Form */
function onCancel() {
  emit('on-cancel')
}
</script>

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
				Товары {{ supplier.name }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					>Отменить</Button
				>
			</div>
		</div>

		<!-- Table -->
		<Card>
			<CardHeader>
				<div class="flex justify-between items-center">
					<div>
						<CardTitle>Материалы поставщика</CardTitle>
						<CardDescription>Цены на материалы поставщика.</CardDescription>
					</div>
				</div>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Товар</TableHead>
							<TableHead>Упаковка</TableHead>
							<TableHead>Цена</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="item in mergedStockMaterialsTableData"
							:key="item.stockMaterialId"
						>
							<TableCell>{{ item.name }}</TableCell>
							<TableCell>{{item.size }} {{ item.unit.name }}</TableCell>
							<TableCell>
								<Input
									v-model.number="item.basePrice"
									type="number"
									class="w-full"
									placeholder="Введите цену"
								/>
							</TableCell>
						</TableRow>
						<TableRow v-if="mergedStockMaterialsTableData.length === 0">
							<TableCell
								colspan="4"
								class="py-4 text-center"
								>Материалы не найдены</TableCell
							>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
		</div>
	</div>
</template>
