<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Switch } from '@/core/components/ui/switch'
import type { IngredientCategoryDTO, IngredientsDTO, UpdateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronLeft } from 'lucide-vue-next'

// Async Components
const AdminSelectIngredientCategory = defineAsyncComponent(() =>
  import('@/modules/admin/ingredient-categories/components/admin-select-ingredient-category.vue'))
const AdminSelectUnit = defineAsyncComponent(() =>
  import('@/modules/admin/units/components/admin-select-unit.vue'))

const { ingredient, readonly = false } = defineProps<{
  ingredient: IngredientsDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateIngredientDTO]
  onCancel: []
}>()

// Validation Schema
const updateIngredientSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название сырья'),
    calories: z.number().min(0, 'Введите корректное значение калорий'),
    fat: z.number().min(0, 'Введите корректное значение жиров'),
    carbs: z.number().min(0, 'Введите корректное значение углеводов'),
    proteins: z.number().min(0, 'Введите корректное значение белков'),
    expirationInDays: z.number().min(0, 'Введите корректные дни хранения'),
    unitId: z.coerce.number().min(1, 'Выберите корректный размер'),
    categoryId: z.coerce.number().min(1, 'Выберите корректную категорию'),
    isAllergen: z.boolean({message: 'Выберите если это сырье является аллергеном'}).default(false)
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: updateIngredientSchema,
  initialValues: {
    name: ingredient.name,
    calories: ingredient.calories,
    fat: ingredient.fat,
    carbs: ingredient.carbs,
    proteins: ingredient.proteins,
    expirationInDays: ingredient.expirationInDays,
    unitId: ingredient.unit.id,
    categoryId: ingredient.category.id,
    isAllergen: ingredient.isAllergen
  }
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (readonly) return

  const dto: UpdateIngredientDTO = {
    name: formValues.name,
    calories: formValues.calories,
    fat: formValues.fat,
    carbs: formValues.carbs,
    proteins: formValues.proteins,
    categoryId: formValues.categoryId,
    unitId: formValues.unitId,
    expirationInDays: formValues.expirationInDays,
    isAllergen: formValues.isAllergen
  }

  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

const openCategoryDialog = ref(false)
const selectedCategory = ref<IngredientCategoryDTO | null>(ingredient.category)

function selectCategory(category: IngredientCategoryDTO) {
  if (readonly) return
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('categoryId', category.id)
}

const openUnitDialog = ref(false)
const selectedUnit = ref<UnitDTO | null>(ingredient.unit)

function selectUnit(unit: UnitDTO) {
  if (readonly) return
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
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
				{{ ingredient.name }}
			</h1>

			<div
				v-if="!readonly"
				class="hidden md:flex items-center gap-2 md:ml-auto"
			>
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
			<!-- LEFT side: Product Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали сырья</CardTitle>
						<CardDescription v-if="!readonly">Заполните информацию об сырьее.</CardDescription>
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
											placeholder="Введите название сырья"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Nutrition Values -->
							<div class="flex gap-4">
								<FormField
									name="calories"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Калории (ккал)</FormLabel>
										<FormControl>
											<Input
												id="calories"
												type="number"
												v-bind="componentField"
												placeholder="Введите калории"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<FormField
									name="fat"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Жиры (грамм)</FormLabel>
										<FormControl>
											<Input
												id="fat"
												type="number"
												v-bind="componentField"
												placeholder="Введите жиры"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>

							<div class="flex gap-4">
								<FormField
									name="carbs"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Углеводы (грамм)</FormLabel>
										<FormControl>
											<Input
												id="carbs"
												type="number"
												v-bind="componentField"
												placeholder="Введите углеводы"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<FormField
									name="proteins"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Белки (грамм)</FormLabel>
										<FormControl>
											<Input
												id="proteins"
												type="number"
												v-bind="componentField"
												placeholder="Введите белки"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>

							<!-- Expiration -->
							<FormField
								name="expirationInDays"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Срок годности (дни)</FormLabel>
									<FormControl>
										<Input
											id="expirationInDays"
											type="number"
											v-bind="componentField"
											placeholder="Введите дни хранения"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ value, handleChange }"
								name="isAllergen"
							>
								<FormItem
									class="flex flex-row justify-between items-center gap-12 p-4 border rounded-lg"
								>
									<div class="flex flex-col space-y-0.5">
										<FormLabel class="font-medium text-base"> Аллерген </FormLabel>
										<FormDescription class="text-sm">
											Данное сырье является аллергеном
										</FormDescription>
									</div>

									<FormControl>
										<Switch
											:checked="value"
											@update:checked="handleChange"
										/>
									</FormControl>
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- RIGHT side: Category & Unit -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Unit Card -->
				<Card>
					<CardHeader>
						<CardTitle>Единица измерения</CardTitle>
						<CardDescription v-if="!readonly">Выберите единицу измерения сырья</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="categoryId">
							<FormItem>
								<template v-if="!readonly">
									<Button
										variant="link"
										class="mt-0 p-0 h-fit text-primary underline"
										@click="openUnitDialog = true"
									>
										{{ selectedUnit?.name || 'Размер не выбран' }}
									</Button>
								</template>
								<template v-else>
									<span
										class="text-muted-foreground"
										>{{ selectedUnit?.name || 'Размер не выбран' }}</span
									>
								</template>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>

				<!-- Category Card -->
				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription v-if="!readonly">Выберите категорию сырья</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="categoryId">
							<FormItem>
								<template v-if="!readonly">
									<Button
										variant="link"
										class="mt-0 p-0 h-fit text-primary underline"
										@click="openCategoryDialog = true"
									>
										{{ selectedCategory?.name || 'Категория не выбрана' }}
									</Button>
								</template>
								<template v-else>
									<span
										class="text-muted-foreground"
										>{{ selectedCategory?.name || 'Категория не выбрана' }}</span
									>
								</template>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Mobile Footer -->
		<div
			v-if="!readonly"
			class="md:hidden flex justify-center items-center gap-2"
		>
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
		<AdminSelectIngredientCategory
			v-if="!readonly"
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>

		<AdminSelectUnit
			v-if="!readonly"
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
	</div>
</template>
