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
				Создать Доставку
			</h1>

			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Supplier Selection -->
		<Card>
			<CardHeader>
				<CardTitle>Выберите поставщика</CardTitle>
				<CardDescription> Укажите поставщика для этой доставки. </CardDescription>
			</CardHeader>
			<CardContent>
				<Button
					variant="link"
					@click="openSupplierDialog = true"
					class="mt-0 p-0 underline"
				>
					{{ selectedSupplier?.name || 'Не выбран' }}
				</Button>
			</CardContent>
		</Card>

		<!-- Materials Table -->
		<Card>
			<CardHeader>
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Материалы доставки</CardTitle>
						<CardDescription class="mt-2">
							Добавьте материалы и укажите их количество и упаковку.
						</CardDescription>
					</div>
					<Button
						variant="outline"
						@click="openStockMaterialDialog = true"
						:disabled="!selectedSupplier"
					>
						{{ !selectedSupplier ? "Выберите поставщика" : "Добавить" }}
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Упаковка</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead>Цена</TableHead>
							<TableHead class="text-center"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow v-if="materials.length === 0">
							<TableCell
								colspan="5"
								class="py-5 text-muted-foreground text-center"
							>
								Нет добавленных материалов
							</TableCell>
						</TableRow>
						<TableRow
							v-for="(material, index) in materials"
							:key="material.stockMaterialId"
						>
							<TableCell>{{ material.name }}</TableCell>
							<TableCell> {{ material.size }} {{ material.unit.name }} </TableCell>
							<TableCell>{{ material.category }}</TableCell>
							<TableCell>
								<Input
									type="number"
									v-model.number="material.quantity"
									:min="1"
									class="w-20"
									:class="{ 'border-red-500': material.quantity <= 0 }"
									placeholder="Введите количество"
								/>
							</TableCell>
							<TableCell>
								<Input
									type="number"
									v-model.number="material.price"
									:min="0"
									class="w-20"
									:class="{ 'border-red-500': material.price <= 0 }"
									placeholder="Введите цену"
								/>
							</TableCell>
							<TableCell class="text-center">
								<Trash
									class="text-red-500 hover:text-red-700 cursor-pointer"
									@click="removeMaterial(index)"
								/>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="onCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>

		<!-- Dialogs -->
		<AdminSelectSupplierDialog
			:open="openSupplierDialog"
			@close="openSupplierDialog = false"
			@select="selectSupplier"
		/>
		<AdminStockMaterialsSelectDialog
			:open="openStockMaterialDialog"
			:initial-filter="stockMaterialFilter"
			@close="openStockMaterialDialog = false"
			@select="addStockMaterial"
		/>
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
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { useToast } from '@/core/components/ui/toast'
import { ChevronLeft, Trash } from 'lucide-vue-next'
import { ref } from 'vue'

// Dialog Components
import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'

// Interfaces
import type { StockMaterialsDTO, StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import AdminSelectSupplierDialog from '@/modules/admin/suppliers/components/admin-select-supplier-dialog.vue'
import type { SupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import type { ReceiveWarehouseDelivery, ReceiveWarehouseStockMaterial } from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'

interface  ReceiveWarehouseStockMaterialForm extends ReceiveWarehouseStockMaterial {
  name: string
  category: string
  size: number
  unit: UnitDTO
}

const emit = defineEmits<{
	(e: 'on-close'): void
	(e: 'on-submit', dto: ReceiveWarehouseDelivery): void
}>()

// Local State
const materials = ref<ReceiveWarehouseStockMaterialForm[]>([])
const selectedSupplier = ref<{ id: number; name: string } | null>(null)

// Dialog States
const openSupplierDialog = ref(false)
const openStockMaterialDialog = ref(false)

const stockMaterialFilter = ref<StockMaterialsFilter>({})

// Toast
const { toast } = useToast()

// Select Supplier
function selectSupplier(supplier: SupplierDTO) {
	selectedSupplier.value = supplier
	openSupplierDialog.value = false
  materials.value = []
  stockMaterialFilter.value = {supplierId: supplier.id}
}

// Add Stock Material
function addStockMaterial(stockMaterial: StockMaterialsDTO) {
	const exists = materials.value.some((item) => item.stockMaterialId === stockMaterial.id)

	if (exists) {
		return
	}

	materials.value.push({
		stockMaterialId: stockMaterial.id,
		quantity: 0,
		price: 0,
		name: stockMaterial.name,
    category: stockMaterial.category.name,
    size: stockMaterial.size,
    unit: stockMaterial.unit
	})

	openStockMaterialDialog.value = false
}

// Remove Material
function removeMaterial(index: number) {
	materials.value.splice(index, 1)
}

// Submit
function onSubmit() {
	if (!selectedSupplier.value) {
		toast({
			title: 'Ошибка',
			description: 'Выберите поставщика.',
			variant: 'destructive',
		})
		return
	}

	if (materials.value.some((item) => item.quantity <= 0)) {
		toast({
			title: 'Ошибка',
			description: 'Убедитесь, что все поля заполнены корректно.',
			variant: 'destructive',
		})
		return
	}

	const payload: ReceiveWarehouseDelivery = {
		supplierId: selectedSupplier.value.id,
		materials: materials.value,
	}

	toast({
		title: 'Успех',
		description: 'Доставка успешно создана.',
		variant: 'default',
	})

	emit('on-submit', payload)
}

// Cancel
function onCancel() {
	emit('on-close')
}
</script>
