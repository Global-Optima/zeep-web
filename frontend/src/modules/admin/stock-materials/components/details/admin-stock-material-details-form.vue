<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

// UI
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import { ChevronLeft, Printer } from 'lucide-vue-next'

// Async Components
const AdminIngredientsSelectDialog = defineAsyncComponent(() =>
  import('@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'))
const AdminSelectStockMaterialCategory = defineAsyncComponent(() =>
  import('@/modules/admin/stock-material-categories/components/admin-select-stock-material-category.vue'))
const AdminSelectUnit = defineAsyncComponent(() =>
  import('@/modules/admin/units/components/admin-select-unit.vue'))

// Types
import { usePrinter } from "@/core/hooks/use-print.hook"
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import type {
  StockMaterialsDTO,
  UpdateStockMaterialDTO,
} from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

// Props
const { stockMaterial, readonly = false } = defineProps<{
  stockMaterial: StockMaterialsDTO
  readonly?: boolean
}>()

// Emit
const emits = defineEmits<{
  onSubmit: [dto: UpdateStockMaterialDTO]
  onCancel: []
}>()

/**
 * Dialog States
 */
const selectedUnit = ref<UnitDTO | null>(stockMaterial.unit || null)
const selectedCategory = ref<StockMaterialCategoryDTO | null>(stockMaterial.category || null)
const selectedIngredient = ref<IngredientsDTO | null>(stockMaterial.ingredient || null)

const openUnitDialog = ref(false)
const openCategoryDialog = ref(false)
const openIngredientDialog = ref(false)

const updateStockMaterialSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название материала'),
    description: z.string().min(1, 'Введите описание'),
    safetyStock: z.coerce.number().min(1, 'Безопасный запас упаковок должен быть больше 0'),
    size: z.coerce.number().min(1, 'Введите размер упаковки'),
    unitId: z.coerce.number().min(1, 'Выберите единицу измерения'),
    categoryId: z.coerce.number().min(1, 'Выберите категорию'),
    ingredientId: z.coerce.number().min(1, 'Выберите ингредиент'),
    expirationPeriodInDays: z.coerce.number().min(1, 'Срок годности должен быть больше 0'),
    barcode: z.string().max(20, 'Введите штрихкод'),
  })
)

const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: updateStockMaterialSchema,
  initialValues: {
    name: stockMaterial.name,
    description: stockMaterial.description,
    safetyStock: stockMaterial.safetyStock,
    unitId: stockMaterial.unit.id,
    categoryId: stockMaterial.category.id,
    ingredientId: stockMaterial.ingredient.id,
    expirationPeriodInDays: stockMaterial.expirationPeriodInDays,
    size: stockMaterial.size,
    barcode: stockMaterial.barcode
  },
})

const onSubmit = handleSubmit((formValues) => {
  if (readonly) return
  const dto: UpdateStockMaterialDTO = formValues
  emits('onSubmit', dto)
})

function onCancel() {
  resetForm()
  emits('onCancel')
}

function selectUnit(unit: UnitDTO) {
  if (readonly) return
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}

function selectCategory(category: StockMaterialCategoryDTO) {
  if (readonly) return
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('categoryId', category.id)
}

function selectIngredient(ingredient: IngredientsDTO) {
  if (readonly) return
  selectedIngredient.value = ingredient
  openIngredientDialog.value = false
  setFieldValue('ingredientId', ingredient.id)
}

const {print} = usePrinter()

const onPrintBarcode = async () => {
  const barcodeBlob = await stockMaterialsService.getBarcodeFile(stockMaterial.id)

  await print(barcodeBlob)
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
				{{ stockMaterial.name }}
			</h1>

			<div
				v-if="!readonly"
				class="md:flex items-center gap-2 hidden md:ml-auto"
			>
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

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<!-- Material Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали материала</CardTitle>
						<CardDescription v-if="!readonly">
							Обновите название, описание и характеристики материала.
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
											:readonly="readonly"
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
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

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
												disabled
											/>

											<Button
												variant="outline"
												@click="onPrintBarcode"
												class="gap-2"
											>
												<Printer class="text-gray-800 size-4" />
												<span>Печать</span>
											</Button>
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
									<FormLabel>Размер упаковки</FormLabel>
									<FormControl>
										<Input
											id="safetyStock"
											type="number"
											v-bind="componentField"
											placeholder="Введите размер упаковки"
											:readonly="readonly"
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
											:readonly="readonly"
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
											:readonly="readonly"
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
				<!-- Unit -->
				<Card>
					<CardHeader>
						<CardTitle>Единица измерения</CardTitle>
						<CardDescription>Выберите единицу измерения.</CardDescription>
					</CardHeader>
					<CardContent>
						<template v-if="!readonly">
							<Button
								variant="link"
								@click="openUnitDialog = true"
								class="mt-0 p-0 underline"
							>
								{{ selectedUnit?.name || 'Не выбрана' }}
							</Button>
						</template>
						<template v-else>
							<span class="text-muted-foreground">{{ selectedUnit?.name || 'Не выбрана' }}</span>
						</template>
					</CardContent>
				</Card>

				<!-- Category -->
				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription>Выберите категорию материала.</CardDescription>
					</CardHeader>
					<CardContent>
						<template v-if="!readonly">
							<Button
								variant="link"
								@click="openCategoryDialog = true"
								class="mt-0 p-0 underline"
							>
								{{ selectedCategory?.name || 'Не выбрана' }}
							</Button>
						</template>
						<template v-else>
							<span
								class="text-muted-foreground"
								>{{ selectedCategory?.name || 'Не выбрана' }}</span
							>
						</template>
					</CardContent>
				</Card>

				<!-- Ingredient -->
				<Card>
					<CardHeader>
						<CardTitle>Ингредиент</CardTitle>
						<CardDescription>Выберите ингредиент.</CardDescription>
					</CardHeader>
					<CardContent>
						<template v-if="!readonly">
							<Button
								variant="link"
								@click="openIngredientDialog = true"
								class="mt-0 p-0 underline"
							>
								{{ selectedIngredient?.name || 'Не выбран' }}
							</Button>
						</template>
						<template v-else>
							<span
								class="text-muted-foreground"
								>{{ selectedIngredient?.name || 'Не выбран' }}</span
							>
						</template>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
		<div
			v-if="!readonly"
			class="flex justify-center items-center gap-2 md:hidden"
		>
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
		<AdminSelectUnit
			v-if="!readonly"
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
		<AdminSelectStockMaterialCategory
			v-if="!readonly"
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>
		<AdminIngredientsSelectDialog
			v-if="!readonly"
			:open="openIngredientDialog"
			@close="openIngredientDialog = false"
			@select="selectIngredient"
		/>
	</div>
</template>
