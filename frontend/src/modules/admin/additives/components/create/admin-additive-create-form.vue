<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import {defineAsyncComponent, ref, useTemplateRef} from 'vue'
import * as z from 'zod'

// UI Components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { Textarea } from '@/core/components/ui/textarea'
import { useToast } from '@/core/components/ui/toast'
import type { AdditiveCategoryDetailsDTO, CreateAdditiveDTO, SelectedIngredientDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { Camera, ChevronLeft, Trash, X } from 'lucide-vue-next'

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

const {isSubmitting} = defineProps<{isSubmitting: boolean}>()

const emits = defineEmits<{
  onSubmit: [dto: CreateAdditiveDTO]
  onCancel: []
}>()

const {toast} = useToast()

// Reactive State for Category Selection
const selectedCategory = ref<AdditiveCategoryDetailsDTO | null>(null)
const openCategoryDialog = ref(false)

const selectedUnit = ref<UnitDTO | null>(null)
const openUnitDialog = ref(false)

const selectedIngredients = ref<SelectedIngredientsTypesDTO[]>([])
const openIngredientsDialog = ref(false)


// Validation Schema
const createAdditiveSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название модификатора'),
    description: z.string().optional(),
    machineId: z.string().min(1, 'Введите код модификатора из автомата').max(40, "Максимум 40 символов"),
    basePrice: z.coerce.number().min(0, 'Введите корректную цену'),
    size: z.coerce.number().min(0, 'Введите размер'),
    unitId: z.coerce.number()
      .min(1, 'Выберите единицу измерения из списка')
      .default(0),
    additiveCategoryId: z.coerce.number()
      .min(1, 'Выберите категорию из списка')
      .default(0),
    image: z.instanceof(File).optional().refine((file) => {
      if (!file) return true; // Optional field
      return ['image/jpeg', 'image/png'].includes(file.type);
    }, 'Поддерживаются только форматы JPEG и PNG').refine((file) => {
      if (!file) return true;
      return file.size <= 5 * 1024 * 1024; // Max 5MB
    }, 'Максимальный размер файла: 5MB'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createAdditiveSchema,
})

const previewImage = ref<string | null>(null);

const imageInputRef = useTemplateRef("imageInputRef");

function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    const file = target.files[0];
    setFieldValue('image', file);
    previewImage.value = URL.createObjectURL(file);
  }
}

function triggerImageInput() {
  imageInputRef.value?.click();
}

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (!selectedCategory.value?.id) return
  if (!selectedUnit.value?.id) return

  if (selectedIngredients.value.some(i => i.quantity <= 0)) {
    return toast({ description: "Укажите количество в технологической карте" })
  }

  const dto: CreateAdditiveDTO = {
    ...formValues,
    additiveCategoryId: selectedCategory.value.id,
    unitId: selectedUnit.value.id,
    ingredients: selectedIngredients.value.map(i => ({ingredientId: i.ingredientId, quantity: i.quantity}))
  }

  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

function selectCategory(category: AdditiveCategoryDetailsDTO) {
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
				:disabled="isSubmitting"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Создать Модификатор
			</h1>

			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					:disabled="isSubmitting"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					:disabled="isSubmitting"
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
						<CardTitle>Детали модификатора</CardTitle>
						<CardDescription>Заполните название, описание и цену модификатора.</CardDescription>
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
											placeholder="Введите название модификатора"
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
											placeholder="Краткое описание модификатора"
											class="min-h-32"
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
												placeholder="Введите цену модификатора"
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
										<FormLabel>
											Размер ({{ selectedUnit?.name.toLowerCase() || 'Не выбрана' }})
										</FormLabel>
										<FormControl>
											<Input
												id="size"
												type="text"
												v-bind="componentField"
												placeholder="500"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>

							<FormField
								name="machineId"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Код модификатора из автомата</FormLabel>
									<FormControl>
										<Input
											id="machineId"
											type="text"
											v-bind="componentField"
											placeholder="Введите код"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>

				<Card class="mt-4">
					<CardHeader>
						<div class="flex justify-between items-start">
							<div>
								<CardTitle>Технологическая карта</CardTitle>
								<CardDescription class="mt-2">
									Выберите инргредиент и его количество
								</CardDescription>
							</div>
							<Button
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
									<TableHead></TableHead>
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
										/>
									</TableCell>
									<TableCell>{{ ingredient.unit }}</TableCell>
									<TableCell class="text-center">
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
						<CardTitle>Изображение</CardTitle>
						<CardDescription>
							Загрузите изображение для продукта.<br />
							Поддерживаемые форматы: JPEG, PNG (макс. 5MB)
						</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="image">
							<FormItem>
								<FormControl>
									<div class="space-y-2">
										<!-- Preview -->
										<div
											v-if="previewImage"
											class="relative w-full h-48"
										>
											<LazyImage
												:src="previewImage"
												alt="Preview"
												class="border rounded-lg w-full h-full object-contain"
											/>
											<button
												type="button"
												class="top-2 right-2 absolute bg-gray-500 transition-all duration-200 hover:bg-red-700 p-1 rounded-full text-white"
												@click="previewImage = null; setFieldValue('image', undefined)"
											>
												<X class="size-4" />
											</button>
										</div>

										<!-- Input -->
										<div
											v-if="!previewImage"
											class="p-4 border-2 border-gray-300 hover:border-primary border-dashed rounded-lg text-center transition-colors cursor-pointer"
											@click="triggerImageInput"
										>
											<input
												ref="imageInputRef"
												type="file"
												accept="image/jpeg, image/png"
												style="display: none;"
												@change="handleImageUpload"
											/>
											<p class="flex flex-col justify-center items-center text-gray-500 text-sm">
												<span class="mb-2"><Camera /></span>
												Нажмите для загрузки изображения<br />
												или перетащите файл
											</p>
										</div>
									</div>
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
						<CardDescription>Выберите категорию продукта</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="additiveCategoryId">
							<FormItem>
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-primary underline"
									@click="openCategoryDialog = true"
								>
									{{ selectedCategory?.name || 'Категория не выбрана' }}
								</Button>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Единица измерения</CardTitle>
						<CardDescription>Выберите единицу измерения</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="unitId">
							<FormItem>
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-primary underline"
									@click="openUnitDialog = true"
								>
									{{ selectedUnit?.name || 'Единица измерения не выбрана' }}
								</Button>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="onCancel"
				:disabled="isSubmitting"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				:disabled="isSubmitting"
				>Сохранить</Button
			>
		</div>

		<AdminIngredientsSelectDialog
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<!-- Category Dialog -->
		<AdminSelectAdditiveCategory
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>

		<AdminSelectUnit
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
	</div>
</template>
