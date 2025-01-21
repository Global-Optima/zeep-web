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

import AdminStockMaterialsSelectDialog from '@/modules/admin/stock-materials/components/admin-stock-materials-select-dialog.vue'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type {
  SupplierDTO,
  SupplierMaterialResponse,
  UpdateSupplierMaterialDTO
} from '@/modules/admin/suppliers/models/suppliers.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

interface UpdateSupplierMaterialForm extends UpdateSupplierMaterialDTO {
  name: string
  unit: UnitDTO
}

/** Props */
const { supplier, stockMaterials } = defineProps<{ supplier: SupplierDTO, stockMaterials: SupplierMaterialResponse[] }>()
const emit = defineEmits<{
  (e: 'on-submit', payload: UpdateSupplierMaterialDTO[]): void
  (e: 'on-cancel'): void
}>()

const { toast } = useToast()

const openDialog = ref(false)

const mergedStockMaterialsTableData = ref<UpdateSupplierMaterialForm[]>(stockMaterials.map(m => ({
    stockMaterialId: m.stockMaterial.id,
    name: m.stockMaterial.name,
    unit: m.stockMaterial.unit,
    basePrice: m.basePrice,
  })))

/** Add New Material */
function addMaterial(material: StockMaterialsDTO) {
  if (mergedStockMaterialsTableData.value.some((item) => item.stockMaterialId === material.id)) {
    toast({
      title: 'Ошибка',
      description: 'Этот материал уже добавлен.',
      variant: 'destructive',
    })
    return
  }

  mergedStockMaterialsTableData.value.push({
    stockMaterialId: material.id,
    name: material.name,
    unit: material.unit,
    basePrice: 0,
  })

  toast({
    title: 'Успех',
    description: `Материал "${material.name}" добавлен.`,
    variant: 'default',
  })
}

/** Remove Material */
function removeMaterial(index: number) {
  const removed = mergedStockMaterialsTableData.value.splice(index, 1)[0]
  toast({
    title: 'Удалено',
    description: `Материал "${removed.name}" удален.`,
    variant: 'default',
  })
}

/** Submit Updated Materials */
function onSubmit() {
  if (mergedStockMaterialsTableData.value.length === 0) {
    toast({
      title: 'Ошибка',
      description: 'Добавьте хотя бы один материал перед отправкой.',
      variant: 'destructive',
    })
    return
  }

  const invalidMaterials = mergedStockMaterialsTableData.value.filter((item) => item.basePrice <= 0)
  if (invalidMaterials.length > 0) {
    toast({
      title: 'Ошибка',
      description: 'Убедитесь, что все цены больше нуля.',
      variant: 'destructive',
    })
    return
  }

  const payload: UpdateSupplierMaterialDTO[] = mergedStockMaterialsTableData.value.map((item) => ({
    stockMaterialId: item.stockMaterialId,
    basePrice: item.basePrice,
  }))

  emit('on-submit', payload)
  toast({
    title: 'Успех',
    description: 'Данные успешно сохранены.',
    variant: 'default',
  })
}

/** Cancel Form */
function onCancel() {
  emit('on-cancel')
  toast({
    title: 'Отмена',
    description: 'Действие отменено.',
    variant: 'default',
  })
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
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Table -->
		<Card>
			<CardHeader>
				<div class="flex justify-between items-center">
					<div>
						<CardTitle>Материалы поставщика</CardTitle>
						<CardDescription>Обновите цены или удалите материалы из списка.</CardDescription>
					</div>
					<Button
						variant="outline"
						@click="openDialog = true"
						>Добавить материал</Button
					>
				</div>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Товар</TableHead>
							<TableHead>Ед. Измерения</TableHead>
							<TableHead>Цена</TableHead>
							<TableHead class="text-center"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="(item, index) in mergedStockMaterialsTableData"
							:key="item.stockMaterialId"
						>
							<TableCell>{{ item.name }}</TableCell>
							<TableCell>{{ item.unit.name }}</TableCell>
							<TableCell>
								<Input
									v-model.number="item.basePrice"
									type="number"
									class="w-full"
									placeholder="Введите цену"
								/>
							</TableCell>
							<TableCell class="text-center">
								<Button
									variant="ghost"
									size="icon"
									@click="removeMaterial(index)"
								>
									<Trash class="text-red-500 hover:text-red-700" />
								</Button>
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
			<Button
				type="submit"
				@click="onSubmit"
				>Сохранить</Button
			>
		</div>

		<!-- Dialog -->
		<AdminStockMaterialsSelectDialog
			:open="openDialog"
			@close="openDialog = false"
			@select="addMaterial"
		/>
	</div>
</template>
