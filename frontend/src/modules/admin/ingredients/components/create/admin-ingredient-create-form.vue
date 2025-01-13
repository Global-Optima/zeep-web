<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import AdminSelectIngredientCategory from '@/modules/admin/ingredient-categories/components/admin-select-ingredient-category.vue'
import type { CreateIngredientDTO, IngredientCategoryDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'

const emits = defineEmits<{
  onSubmit: [dto: CreateIngredientDTO]
  onCancel: []
}>()

const { data: units } = useQuery({
  queryKey:  ['admin-units'],
	queryFn: () => unitsService.getAllUnits(),
})

const { data: ingredientCategories } = useQuery({
  queryKey:  ['admin-ingredient-categories'],
	queryFn: () => ingredientsService.getIngredientCategories(),
})

// Validation Schema
const updateIngredientSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название ингредиента'),
    calories: z.number().min(0, 'Введите корректное значение калорий'),
    fat: z.number().min(0, 'Введите корректное значение жиров'),
    carbs: z.number().min(0, 'Введите корректное значение углеводов'),
    proteins: z.number().min(0, 'Введите корректное значение белков'),
    expirationInDays: z.number().min(0, 'Введите корректные дни хранения'),
    unitId: z.coerce.number().min(1, 'Выберите корректный размер'),
    categoryId: z.coerce.number().min(1, 'Выберите корректную категорию')
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: updateIngredientSchema,
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  const dto: CreateIngredientDTO = {
    name: formValues.name,
    calories: formValues.calories,
    fat: formValues.fat,
    carbs: formValues.carbs,
    proteins: formValues.proteins,
    categoryId: formValues.categoryId,
    unitId: formValues.unitId,
    expirationInDays: formValues.expirationInDays
  }

  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

const openCategoryDialog = ref(false)
const selectedCategory = ref<IngredientCategoryDTO | null>(null)

function selectCategory(category: IngredientCategoryDTO) {
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('categoryId', category.id)
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
				Создать ингредиент
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
			<!-- LEFT side: Product Details (Name, Description) -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали ингредиента</CardTitle>
						<CardDescription>Заполните информацию об ингредиенте.</CardDescription>
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
											placeholder="Введите название ингредиента"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Calories, Fat, Carbs, Proteins -->
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
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>

							<!-- Expiration Date -->
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
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- RIGHT side: Media & Category -->
			<div class="items-start gap-4 grid auto-rows-max">
				<Card>
					<!-- Category Card -->
					<CardHeader>
						<CardTitle>Размер</CardTitle>
						<CardDescription>Выберите размер ингредиента</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField
							name="unitId"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormControl>
									<Select v-bind="componentField">
										<SelectTrigger class="w-full">
											<SelectValue
												class="text-sm sm:text-base"
												placeholder="Выберите размер"
											/>
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="unit in units"
												:key="unit.id"
												:value="unit.id.toString()"
											>
												{{ unit.name }}
											</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>

				<Card>
					<!-- Category Card -->
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription>Выберите категорию ингредиента</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="categoryId">
							<FormItem>
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-blue-600 underline"
									@click="openCategoryDialog = true"
								>
									{{ selectedCategory?.name || 'Категория не выбрана' }}
								</Button>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>
			</div>
		</div>

		<AdminSelectIngredientCategory
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>

		<!-- Footer -->
		<div class="flex justify-center items-center gap-2 md:hidden">
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
	</div>
</template>
