<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { Textarea } from '@/core/components/ui/textarea'
import { useToast } from '@/core/components/ui/toast'
import type { AdditiveCategoryDTO, AdditiveDetailsDTO, BaseAdditiveCategoryDTO, SelectedIngredientDTO, UpdateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'

// Async Components
const AdminSelectAdditiveCategory = defineAsyncComponent(() =>
  import('@/modules/admin/additive-categories/components/admin-select-additive-category.vue'))
const AdminIngredientsSelectDialog = defineAsyncComponent(() =>
  import('@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'))
const AdminSelectUnit = defineAsyncComponent(() =>
  import('@/modules/admin/units/components/admin-select-unit.vue'))

interface SelectedIngredientsTypesDTO extends SelectedIngredientDTO {
  name: string
  unit: string
  category: string
}

const { additive, readonly = false } = defineProps<{
  additive: AdditiveDetailsDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateAdditiveDTO]
  onCancel: []
}>()

const { toast } = useToast()

// Reactive State
const selectedCategory = ref<BaseAdditiveCategoryDTO | null>(additive.category)
const openCategoryDialog = ref(false)
const selectedUnit = ref<UnitDTO | null>(additive.unit)
const openUnitDialog = ref(false)
const selectedIngredients = ref<SelectedIngredientsTypesDTO[]>(additive.ingredients.map(i => ({
  ingredientId: i.ingredient.id,
  quantity: i.quantity,
  name: i.ingredient.name,
  unit: i.ingredient.unit.name,
  category: i.ingredient.category.name,
})))
const openIngredientsDialog = ref(false)

// Validation Schema
const createAdditiveSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название добавки'),
    description: z.string().min(1, 'Введите описание'),
    basePrice: z.coerce.number().min(0, 'Введите корректную цену'),
    size: z.coerce.number().min(0, 'Введите размер'),
    unitId: z.number().min(0, 'Введите единицу измерения'),
    imageUrl: z.string().min(1, 'Вставьте картинку добавки'),
    additiveCategoryId: z.coerce.number().min(1, 'Выберите категорию добавки'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createAdditiveSchema,
  initialValues: {
    name: additive.name,
    description: additive.description,
    basePrice: additive.basePrice,
    size: additive.size,
    unitId: additive.unit.id,
    imageUrl: additive.imageUrl,
    additiveCategoryId: additive.category.id,
  }
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (readonly) return

  if (!selectedCategory.value?.id || !selectedUnit.value?.id) return
  if (selectedIngredients.value.length === 0) {
    return toast({ description: "Технологическая карта должна иметь минимум 1 ингредиент" })
  }
  if (selectedIngredients.value.some(i => i.quantity <= 0)) {
    return toast({ description: "Технологическая карта не может иметь количество 0" })
  }

  const dto: UpdateAdditiveDTO = {
    ...formValues,
    additiveCategoryId: selectedCategory.value.id,
    unitId: selectedUnit.value.id,
    ingredients: selectedIngredients.value.map(i => ({ ingredientId: i.ingredientId, quantity: i.quantity }))
  }

  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

function selectCategory(category: AdditiveCategoryDTO) {
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('additiveCategoryId', category.id)
}

function selectUnit(unit: UnitDTO) {
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}

function addIngredient(ingredient: IngredientsDTO) {
  if (!selectedIngredients.value.some((item) => item.ingredientId === ingredient.id)) {
    selectedIngredients.value.push({
      ingredientId: ingredient.id,
      name: ingredient.name,
      unit: ingredient.unit.name,
      category: ingredient.category.name,
      quantity: 0
    })
  }
}
function removeIngredient(index: number) {
  selectedIngredients.value.splice(index, 1)
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
				{{ additive.name }}
			</h1>

			<div
				v-if="!readonly"
				class="md:flex items-center gap-2 hidden md:ml-auto"
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
			<!-- Additive Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали добавки</CardTitle>
						<CardDescription v-if="!readonly"
							>Заполните название, описание и цену добавки.</CardDescription
						>
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
											placeholder="Введите название добавки"
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
											placeholder="Краткое описание добавки"
											class="min-h-32"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Price and Size -->
							<div class="flex gap-4">
								<FormField
									name="basePrice"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Цена</FormLabel>
										<FormControl>
											<Input
												id="price"
												type="number"
												v-bind="componentField"
												placeholder="Введите цену добавки"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>

								<FormField
									name="size"
									v-slot="{ componentField }"
								>
									<FormItem class="flex-1">
										<FormLabel>Размер</FormLabel>
										<FormControl>
											<Input
												id="size"
												type="text"
												v-bind="componentField"
												placeholder="500"
												:readonly="readonly"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>
						</div>
					</CardContent>
				</Card>

				<Card class="mt-4">
					<CardHeader>
						<div class="flex justify-between items-start">
							<div>
								<CardTitle>Технологическая карта</CardTitle>
								<CardDescription
									v-if="!readonly"
									class="mt-2"
								>
									Выберите инргредиент и его количество
								</CardDescription>
							</div>
							<Button
								v-if="!readonly"
								variant="outline"
								@click="openIngredientsDialog = true"
							>
								Добавить
							</Button>
						</div>
					</CardHeader>
					<CardContent>
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Название</TableHead>
									<TableHead>Категория</TableHead>
									<TableHead>Количество</TableHead>
									<TableHead>Размер</TableHead>
									<TableHead v-if="!readonly"></TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								<TableRow
									v-for="(ingredient, index) in selectedIngredients"
									:key="ingredient.ingredientId"
								>
									<TableCell>{{ ingredient.name }}</TableCell>
									<TableCell>{{ ingredient.category }}</TableCell>
									<TableCell class="flex items-center gap-2">
										<Input
											type="number"
											v-model.number="ingredient.quantity"
											:min="0"
											placeholder="Введите количество"
											class="w-16"
											:readonly="readonly"
										/>
									</TableCell>
									<TableCell>{{ ingredient.unit }}</TableCell>
									<TableCell
										v-if="!readonly"
										class="text-center"
									>
										<Button
											variant="ghost"
											size="icon"
											@click="removeIngredient(index)"
										>
											<Trash class="w-6 h-6 text-red-500" />
										</Button>
									</TableCell>
								</TableRow>
							</TableBody>
						</Table>
					</CardContent>
				</Card>
			</div>

			<!-- Media and Category Blocks -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Media Block -->
				<Card>
					<CardHeader>
						<CardTitle>Медиа</CardTitle>
						<CardDescription v-if="!readonly">Загрузите изображение добавки.</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField
							name="imageUrl"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Изображение</FormLabel>
								<FormControl>
									<Input
										id="imageUrl"
										type="text"
										v-bind="componentField"
										placeholder="Вставьте ссылку на изображение"
										:readonly="readonly"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>

				<!-- Category Block -->
				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription v-if="!readonly">Выберите категорию топпинга.</CardDescription>
					</CardHeader>
					<CardContent>
						<div>
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
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Единица измерения</CardTitle>
						<CardDescription v-if="!readonly">Выберите единицу измерения</CardDescription>
					</CardHeader>
					<CardContent>
						<div>
							<template v-if="!readonly">
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-primary underline"
									@click="openUnitDialog = true"
								>
									{{ selectedUnit?.name || 'Единица измерения не выбрана' }}
								</Button>
							</template>
							<template v-else>
								<span
									class="text-muted-foreground"
									>{{ selectedUnit?.name || 'Единица измерения не выбрана' }}</span
								>
							</template>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Mobile Footer -->
		<div
			v-if="!readonly"
			class="flex justify-center items-center gap-2 md:hidden"
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
		<AdminIngredientsSelectDialog
			v-if="!readonly"
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<AdminSelectAdditiveCategory
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
