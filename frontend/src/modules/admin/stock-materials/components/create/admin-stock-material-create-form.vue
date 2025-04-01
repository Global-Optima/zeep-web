<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'

import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import AdminSelectStockMaterialCategory from '@/modules/admin/stock-material-categories/components/admin-select-stock-material-category.vue'
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'

import { ChevronLeft } from 'lucide-vue-next'

import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import type { CreateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

const emits = defineEmits<{
  onSubmit: [dto: CreateStockMaterialDTO]
  onCancel: []
}>()


const selectedUnit = ref<UnitDTO | null>(null)
const selectedCategory = ref<StockMaterialCategoryDTO | null>(null)
const selectedIngredient = ref<IngredientsDTO | null>(null)

const openUnitDialog = ref(false)
const openCategoryDialog = ref(false)
const openIngredientDialog = ref(false)

const createStockMaterialSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название материала'),
    description: z.string().optional(),
    safetyStock: z.coerce.number().min(1, 'Безопасный запас упаковок должен быть больше 0'),
    size: z.coerce.number().min(1, 'Введите размер упаковки'),
    unitId: z.coerce.number().min(1, 'Выберите единицу измерения'),
    categoryId: z.coerce.number().min(1, 'Выберите категорию'),
    ingredientId: z.coerce.number().min(1, 'Выберите ингредиент'),
    barcode: z
      .string()
      .min(1, 'Создайте или вставьте штрихкод')
      .max(20, 'Введите штрихкод'),
    expirationPeriodInDays: z.coerce.number().min(1, 'Срок годности должен быть больше 0'),
  })
)

const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createStockMaterialSchema,
})

const onSubmit = handleSubmit((formValues) => {
  const dto: CreateStockMaterialDTO = {
    ...formValues,
    isActive: true
  }
  emits('onSubmit', dto)
  resetForm()
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
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

const onGenerateBarcodeClick = async () => {
  const barcodeResponse = await stockMaterialsService.generateBarcode()
  setFieldValue('barcode', barcodeResponse.barcode)
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
				Создать Материал
			</h1>

			<div class="hidden md:flex items-center gap-2 md:ml-auto">
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
						<CardDescription>
							Заполните название, описание и характеристики материала.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<!-- Name -->
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

							<!-- Description -->
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

							<!-- Barcode -->
							<FormField
								name="barcode"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Штрихкод</FormLabel>
									<FormControl>
										<div class="flex items-center gap-2">
											<Input
												id="barcode"
												v-bind="componentField"
												placeholder="Введите штрихкод"
											/>

											<Button
												variant="outline"
												@click="onGenerateBarcodeClick"
												>Создать</Button
											>
										</div>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="size"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel
										>Размер упаковки ({{ selectedUnit?.name.toLowerCase() || 'Не выбрана' }})
									</FormLabel>
									<FormControl>
										<Input
											id="safetyStock"
											type="number"
											v-bind="componentField"
											placeholder="Введите размер упаковки"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Safety Stock -->
							<FormField
								name="safetyStock"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Безопасный запас упаковок</FormLabel>
									<FormControl>
										<Input
											id="safetyStock"
											type="number"
											v-bind="componentField"
											placeholder="Введите безопасный запас упаковок"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Expiration Period -->
							<FormField
								name="expirationPeriodInDays"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Срок годности (дни)</FormLabel>
									<FormControl>
										<Input
											id="expirationPeriodInDays"
											type="number"
											v-bind="componentField"
											placeholder="Введите минимальный срок годности в днях"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Right Column: Unit, Category, Ingredient -->
			<div class="items-start gap-4 grid auto-rows-max">
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

		<!-- Footer for mobile -->
		<div class="md:hidden flex justify-center items-center gap-2">
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

		<!-- Dialogs -->
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
