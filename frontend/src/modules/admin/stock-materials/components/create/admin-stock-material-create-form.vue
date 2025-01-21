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
import { useToast } from '@/core/components/ui/toast'

// Select Dialog Components
import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import AdminSelectStockMaterialCategory from '@/modules/admin/stock-material-categories/components/admin-select-stock-material-category.vue'
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'

import { ChevronLeft, Trash } from 'lucide-vue-next'

// Typings
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import type { CreateStockMaterialDTO, CreateStockMaterialPackagesDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

/**
 * Emits
 */
const emits = defineEmits<{
  onSubmit: [dto: CreateStockMaterialDTO]
  onCancel: []
}>()

/**
 * Dialog States
 */
const selectedUnit = ref<UnitDTO | null>(null)
const selectedCategory = ref<StockMaterialCategoryDTO | null>(null)
const selectedIngredient = ref<IngredientsDTO | null>(null)

const openUnitDialog = ref(false)
const openCategoryDialog = ref(false)
const openIngredientDialog = ref(false)

/**
 * For Package Rows
 */
interface PackageRow {
  size: number
  unitId: number
  unitName?: string
}
const packages = ref<PackageRow[]>([])
const openPackageUnitDialog = ref(false)
/** This holds the index of the package row we’re currently editing. */
const selectedPackageIndex = ref<number | null>(null)

/** Toast */
const { toast } = useToast()

/**
 * Validation Schema
 */
const createStockMaterialSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название материала'),
    description: z.string().optional(),
    safetyStock: z.coerce.number().min(1, 'Безопасный запас должен быть больше 0'),
    unitId: z.coerce.number().min(1, 'Выберите единицу измерения'),
    categoryId: z.coerce.number().min(1, 'Выберите категорию'),
    ingredientId: z.coerce.number().min(1, 'Выберите ингредиент'),
    barcode: z
      .string()
      .min(1, 'Введите штрихкод')
      .length(12, 'Штрихкод должен быть длинной в 12 символов'),
    expirationPeriodInDays: z.coerce.number().min(1, 'Срок годности должен быть больше 0'),
  })
)

/**
 * Setup the form
 */
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createStockMaterialSchema,
})

/**
 * Submissions
 */
const onSubmit = handleSubmit((formValues) => {
  // 1) Validate “packages” array
  if (packages.value.length === 0) {
    toast({
      title: 'Ошибка',
      description: 'Добавьте хотя бы одну упаковку',
      variant: 'destructive',
    })
    return
  }

  for (const pkg of packages.value) {
    if (pkg.size <= 0 || pkg.unitId <= 0) {
      toast({
        title: 'Ошибка',
        description: 'Все упаковки должны иметь корректный размер и единицу измерения.',
        variant: 'destructive',
      })
      return
    }
  }

  // 2) Build final payload with packages
  const finalPackages: CreateStockMaterialPackagesDTO[] = packages.value.map((p) => ({
    size: p.size,
    unitId: p.unitId,
  }))

  const dto: CreateStockMaterialDTO = {
    ...formValues,
    packages: finalPackages,
  }

  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  packages.value = []
  emits('onCancel')
}

/**
 * Selection Handlers
 */
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

/**
 * Packages Methods
 */
function addPackageRow() {
  packages.value.push({
    size: 0,
    unitId: 0,
    unitName: 'Не выбрана',
  })
}
function removePackageRow(index: number) {
  packages.value.splice(index, 1)
}
function openPackageRowUnitDialog(index: number) {
  selectedPackageIndex.value = index
  openPackageUnitDialog.value = true
}
/** Called after user picks a unit for a specific package row */
function selectPackageUnit(unit: UnitDTO) {
  if (selectedPackageIndex.value === null) return
  const row = packages.value[selectedPackageIndex.value]
  row.unitId = unit.id
  row.unitName = unit.name
  openPackageUnitDialog.value = false
  selectedPackageIndex.value = null
  toast({
    title: 'Упаковка',
    description: `Выбрана единица измерения: ${unit.name}`,
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
				Создать Материал
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
										<Input
											id="barcode"
											v-bind="componentField"
											placeholder="Введите штрихкод"
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

				<!-- Packages Card -->
				<Card>
					<CardHeader>
						<div class="flex justify-between items-start gap-4">
							<div>
								<CardTitle>Упаковки</CardTitle>
								<CardDescription class="mt-2">Укажите возможные варианты упаковки.</CardDescription>
							</div>

							<Button
								variant="outline"
								@click="addPackageRow"
							>
								Добавить
							</Button>
						</div>
					</CardHeader>
					<CardContent>
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Размер</TableHead>
									<TableHead>Ед. Измерения</TableHead>
									<TableHead class="text-center"></TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								<!-- Render each package row -->
								<TableRow
									v-for="(pkg, index) in packages"
									:key="index"
								>
									<!-- Size -->
									<TableCell>
										<Input
											type="number"
											class="w-24"
											v-model.number="pkg.size"
											placeholder="Размер"
										/>
									</TableCell>

									<!-- Unit -->
									<TableCell>
										<span
											class="text-primary underline cursor-pointer"
											@click="openPackageRowUnitDialog(index)"
										>
											{{ pkg.unitName || 'Не выбрана' }}
										</span>
									</TableCell>
									<!-- Actions -->
									<TableCell class="text-center">
										<Button
											variant="ghost"
											size="icon"
											@click="removePackageRow(index)"
										>
											<Trash class="w-5 h-5 text-destructive" />
										</Button>
									</TableCell>
								</TableRow>

								<!-- Empty state -->
								<TableRow v-if="packages.length === 0">
									<TableCell
										colspan="3"
										class="py-2 text-center text-muted-foreground"
									>
										Пакетов нет
									</TableCell>
								</TableRow>
							</TableBody>
						</Table>
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
		<AdminSelectUnit
			:open="openPackageUnitDialog"
			@close="openPackageUnitDialog = false"
			@select="selectPackageUnit"
		/>
	</div>
</template>
