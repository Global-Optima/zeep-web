<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'

// Select Dialog Components
import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import AdminSelectStockMaterialCategory from '@/modules/admin/stock-material-categories/components/admin-select-stock-material-category.vue'
import AdminSelectSupplierDialog from '@/modules/admin/suppliers/components/admin-select-supplier-dialog.vue'
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'

// Typings
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import type { StockMaterialsDTO, UpdateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { SuppliersDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

import { ChevronLeft } from 'lucide-vue-next'

// Props
const {stockMaterial} = defineProps<{
  stockMaterial: StockMaterialsDTO
}>()

// Emit Events
const emits = defineEmits<{
  onSubmit: [dto: UpdateStockMaterialDTO]
  onCancel: []
}>()

// Dialog States
const selectedSupplier = ref<SuppliersDTO | null>(stockMaterial.supplier || null)
const selectedUnit = ref<UnitDTO | null>(stockMaterial.unit || null)
const selectedCategory = ref<StockMaterialCategoryDTO | null>(stockMaterial.category || null)
const selectedIngredient = ref<IngredientsDTO | null>(stockMaterial.ingredient || null)

const openSupplierDialog = ref(false)
const openUnitDialog = ref(false)
const openCategoryDialog = ref(false)
const openIngredientDialog = ref(false)

// Validation Schema
const updateStockMaterialSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название материала'),
    description: z.string().optional(),
    safetyStock: z.coerce.number().min(1, 'Безопасный запас должен быть больше 0'),
    unitId: z.coerce.number().min(1, 'Выберите единицу измерения'),
    supplierId: z.coerce.number().min(1, 'Выберите поставщика'),
    categoryId: z.coerce.number().min(1, 'Выберите категорию'),
    ingredientId: z.coerce.number().min(1, 'Выберите ингредиент'),
    barcode: z.string().optional(),
    expirationPeriodInDays: z.coerce.number().min(1, 'Срок годности должен быть больше 0'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm<UpdateStockMaterialDTO>({
  validationSchema: updateStockMaterialSchema,
  initialValues: {
    name: stockMaterial.name,
    description: stockMaterial.description,
    safetyStock: stockMaterial.safetyStock,
    unitId: stockMaterial.unit.id,
    supplierId: stockMaterial?.supplier?.id ?? 1, // TODO: add proper supplier id
    categoryId: stockMaterial.category.id,
    ingredientId: stockMaterial.ingredient.id,
    barcode: stockMaterial.barcode,
    expirationPeriodInDays: stockMaterial.expirationPeriodInDays,
  },
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  emits('onSubmit', formValues)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

// Selection Handlers
function selectSupplier(supplier: SuppliersDTO) {
  selectedSupplier.value = supplier
  openSupplierDialog.value = false
  setFieldValue('supplierId', supplier.id)
}

function selectUnit(unit: UnitDTO) {
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}

function selectCategory(category: StockMaterialCategoryDTO) {
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('categoryId', category.id)
}

function selectIngredient(ingredient: IngredientsDTO) {
  selectedIngredient.value = ingredient
  openIngredientDialog.value = false
  setFieldValue('ingredientId', ingredient.id)
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
				Обновить Материал
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

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<!-- Material Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали материала</CardTitle>
						<CardDescription
							>Обновите название, описание и характеристики материала.</CardDescription
						>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<FormField
								name="name"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Название</FormLabel>
									<FormControl>
										<Input
											id="name"
											type="text"
											v-bind="componentField"
											placeholder="Введите название материала"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="description"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Описание</FormLabel>
									<FormControl>
										<Textarea
											id="description"
											v-bind="componentField"
											placeholder="Краткое описание материала"
											class="min-h-32"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="safetyStock"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Безопасный запас</FormLabel>
									<FormControl>
										<Input
											id="safetyStock"
											type="number"
											v-bind="componentField"
											placeholder="Введите безопасный запас"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Select Dialogs -->
			<div class="items-start gap-4 grid auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Поставщик</CardTitle>
						<CardDescription>Выберите поставщика.</CardDescription>
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

				<Card>
					<CardHeader>
						<CardTitle>Единица измерения</CardTitle>
						<CardDescription>Выберите единицу измерения.</CardDescription>
					</CardHeader>
					<CardContent>
						<Button
							variant="link"
							@click="openUnitDialog = true"
							class="mt-0 p-0 underline"
						>
							{{ selectedUnit?.name || 'Не выбрана' }}
						</Button>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription>Выберите категорию материала.</CardDescription>
					</CardHeader>
					<CardContent>
						<Button
							variant="link"
							@click="openCategoryDialog = true"
							class="mt-0 p-0 underline"
						>
							{{ selectedCategory?.name || 'Не выбрана' }}
						</Button>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Ингредиент</CardTitle>
						<CardDescription>Выберите ингредиент.</CardDescription>
					</CardHeader>
					<CardContent>
						<Button
							variant="link"
							@click="openIngredientDialog = true"
							class="mt-0 p-0 underline"
						>
							{{ selectedIngredient?.name || 'Не выбран' }}
						</Button>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
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

		<!-- Select Dialog Components -->
		<AdminSelectSupplierDialog
			:open="openSupplierDialog"
			@close="openSupplierDialog = false"
			@select="selectSupplier"
		/>
		<AdminSelectUnit
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
		<AdminSelectStockMaterialCategory
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>
		<AdminIngredientsSelectDialog
			:open="openIngredientDialog"
			@close="openIngredientDialog = false"
			@select="selectIngredient"
		/>
	</div>
</template>
