<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref, useTemplateRef } from 'vue'
import * as z from 'zod'

// UI Components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/core/components/ui/dropdown-menu'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { Textarea } from '@/core/components/ui/textarea'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from "@/core/config/routes.config"
import { TechnicalMapEntity, useCopyTechnicalMap } from '@/core/hooks/use-copy-technical-map.hooks'
import type { AdditiveCategoryDetailsDTO, AdditiveCategoryDTO, AdditiveDetailsDTO, SelectedIngredientDTO, UpdateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { ProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { Camera, ChevronLeft, EllipsisVertical, Trash, X } from 'lucide-vue-next'

// Async Components
const AdminSelectAdditiveCategory = defineAsyncComponent(() =>
  import('@/modules/admin/additive-categories/components/admin-select-additive-category.vue'))
const AdminIngredientsSelectDialog = defineAsyncComponent(() =>
  import('@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'))
const AdminSelectUnit = defineAsyncComponent(() =>
  import('@/modules/admin/units/components/admin-select-unit.vue'))
const AdminSelectProvisionDialog = defineAsyncComponent(() =>
  import("@/modules/admin/provisions/components/admin-select-provision-dialog.vue"))

interface SelectedIngredientsTypesDTO extends SelectedIngredientDTO {
  name: string
  unit: string
  category: string
}

interface SelectedProvisionsTypesDTO {
  provisionId: number
  name: string
  absoluteVolume: number
  unit: string
  volume: number
}

const { additive, readonly = false, isSubmitting } = defineProps<{
  additive: AdditiveDetailsDTO
  readonly?: boolean
  isSubmitting: boolean
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateAdditiveDTO]
  onCancel: []
}>()

const { toast } = useToast()
const { setTechnicalMapReference, fetchTechnicalMap } = useCopyTechnicalMap()


// Reactive State
const selectedCategory = ref<AdditiveCategoryDTO | null>(additive.category)
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
const updateAdditiveSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название модификатора')
      .max(100, 'Название не может превышать 100 символов'),
    description: z.string()
      .max(500, 'Описание не может превышать 500 символов').optional(),
    machineId: z.string().min(1, 'Введите код модификатора из автомата').max(40, "Максимум 40 символов"),
    basePrice: z.coerce.number().min(0, 'Введите корректную цену'),
    size: z.coerce.number().min(0, 'Введите размер'),
    unitId: z.number().min(0, 'Введите единицу измерения'),
    additiveCategoryId: z.coerce.number().min(1, 'Выберите категорию модификатора'),
    image: z.instanceof(File).optional().refine((file) => {
      if (!file) return true;
      return ['image/jpeg', 'image/png'].includes(file.type);
    }, 'Поддерживаются только форматы JPEG и PNG').refine((file) => {
      if (!file) return true;
      return file.size <= 5 * 1024 * 1024;
    }, 'Максимальный размер файла: 5MB'),
  })
)

// Form Setup
const { handleSubmit, resetField, setFieldValue } = useForm({
  validationSchema: updateAdditiveSchema,
  initialValues: {
    name: additive.name,
    description: additive.description,
    basePrice: additive.basePrice,
    size: additive.size,
    unitId: additive.unit.id,
    additiveCategoryId: additive.category.id,
    machineId: additive.machineId,
  }
})

const previewImage = ref<string | null>(additive.imageUrl || null);
const deleteImage = ref<boolean>(false)

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.length) {
    const file = target.files[0];
    previewImage.value = URL.createObjectURL(file);
    setFieldValue('image', file);
    deleteImage.value = true
  }
};

const provisions = ref<SelectedProvisionsTypesDTO[]>(additive.provisions.map(p => ({
  provisionId: p.provision.id,
  name: p.provision.name,
  absoluteVolume: p.provision.absoluteVolume,
  unit: p.provision.unit.name,
  volume: p.volume,
})))
const openProvisionsDialog = ref(false)

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (readonly) return

  if (!selectedCategory.value?.id || !selectedUnit.value?.id) return

  if (selectedIngredients.value.some(i => i.quantity <= 0)) {
    return toast({ description: "Укажите количество в технологической карте" })
  }

  if (provisions.value.some(i => i.volume <= 0)) {
    return toast({ description: "Укажите количество в технологической карте" })
  }

  const dto: UpdateAdditiveDTO = {
    ...formValues,
    additiveCategoryId: selectedCategory.value.id,
    unitId: selectedUnit.value.id,
    ingredients: selectedIngredients.value.map(i => ({ ingredientId: i.ingredientId, quantity: i.quantity })),
    deleteImage: deleteImage.value,
    provisions: provisions.value.map(p => ({provisionId: p.provisionId, volume: p.volume}))
  }

  emits('onSubmit', dto)
})

const resetFormValues = () => {
  resetField('image')
  deleteImage.value = false
}

defineExpose({ resetFormValues })

const onDeleteImage = () => {
  previewImage.value = null
  setFieldValue('image', undefined)
  deleteImage.value = true
}

const onCancel = () => {
  emits('onCancel')
}

function addProvision(provision: ProvisionDTO) {
  if (!provisions.value.some((item) => item.provisionId === provision.id)) {
    provisions.value.push({
      provisionId: provision.id,
      name: provision.name,
      unit: provision.unit.name,
      absoluteVolume: provision.absoluteVolume,
      volume: 0
    })
  }
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

function removeProvision(index: number) {
  if (readonly) return
  provisions.value.splice(index, 1)
}

const imageInputRef = useTemplateRef("imageInputRef");
const triggerImageInput = () => imageInputRef.value?.click();

const onCopyTechMapClick = () => {
  try {
    setTechnicalMapReference(TechnicalMapEntity.ADDITIVE, additive.id)
    toast({ description: "Технологическая карта успешно скопирована", variant: "success" })

  } catch {
    toast({ description: "Ошибка при копировании технологической карты", variant: "destructive" })
  }
}

const onPasteTechMapClick = async () => {
  try {
    const techMap = await fetchTechnicalMap()
    if (!techMap) {
      toast({ description: "Технологическая карта не найдена" })
      return
    }

    selectedIngredients.value = techMap.map(t => ({
      ingredientId: t.ingredient.id,
      name: t.ingredient.name,
      unit: t.ingredient.unit.name,
      category: t.ingredient.category.name,
      quantity: t.quantity,
    }))

    toast({ description: "Технологическая карта успешно вставлена", variant: "success" })

  } catch {
    toast({ description: "Ошибка при вставке технологической карты", variant: "destructive" })
  }
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
				{{ additive.name }}
			</h1>

			<div
				v-if="!readonly"
				class="hidden md:flex items-center gap-2 md:ml-auto"
			>
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
						<CardDescription v-if="!readonly"
							>Заполните название, описание и цену модификатора.</CardDescription
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
											placeholder="Введите название модификатора"
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
											placeholder="Краткое описание модификатора"
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
												placeholder="Введите цену модификатора"
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
										<FormLabel>
											Размер ({{ selectedUnit?.name.toLowerCase() || 'Не выбрана' }})
										</FormLabel>
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
								<CardDescription
									v-if="!readonly"
									class="mt-2"
								>
									Выберите инргредиент и его количество
								</CardDescription>
							</div>
							<div
								v-if="!readonly"
								class="flex items-center gap-2"
							>
								<DropdownMenu>
									<DropdownMenuTrigger class="p-2 border rounded-md">
										<EllipsisVertical class="size-4" />
									</DropdownMenuTrigger>
									<DropdownMenuContent>
										<DropdownMenuItem @click="onCopyTechMapClick">Копировать</DropdownMenuItem>
										<DropdownMenuItem @click="onPasteTechMapClick">Вставить</DropdownMenuItem>
									</DropdownMenuContent>
								</DropdownMenu>

								<Button
									variant="outline"
									@click="openIngredientsDialog = true"
								>
									Добавить
								</Button>
							</div>
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
									<TableCell>
										<RouterLink
											:to="{name: getRouteName('ADMIN_INGREDIENTS_DETAILS'), params: {id: ingredient.ingredientId}}"
											target="_blank"
											class="hover:text-primary underline underline-offset-4 transition-colors duration-300"
										>
											{{ ingredient.name }}
										</RouterLink>
									</TableCell>
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

				<Card class="mt-4">
					<CardHeader>
						<div class="flex justify-between items-start">
							<div>
								<CardTitle>Заготовки</CardTitle>
								<CardDescription
									v-if="!readonly"
									class="mt-2"
								>
									Выберите заготовки и их обьем
								</CardDescription>
							</div>
							<Button
								v-if="!readonly"
								variant="outline"
								@click="openProvisionsDialog = true"
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
									<TableHead>Изначальный обьем</TableHead>
									<TableHead>Обьем для продукта</TableHead>
									<TableHead v-if="!readonly"></TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								<TableRow
									v-for="(provision, index) in provisions"
									:key="provision.provisionId"
								>
									<TableCell>
										<RouterLink
											:to="{name: getRouteName('ADMIN_PROVISION_DETAILS'), params: {id: provision.provisionId}}"
											target="_blank"
											class="hover:text-primary underline underline-offset-4 transition-colors duration-300"
										>
											{{ provision.name }}
										</RouterLink>
									</TableCell>
									<TableCell
										>{{ provision.absoluteVolume }} {{ provision.unit.toLowerCase() }}</TableCell
									>

									<TableCell class="flex items-center gap-4">
										<Input
											type="number"
											v-model.number="provision.volume"
											:min="0"
											class="w-24"
											placeholder="Введите нужный обьем"
										/>
										{{ provision.unit.toLowerCase() }}
									</TableCell>
									<TableCell
										v-if="!readonly"
										class="text-center"
									>
										<Button
											variant="ghost"
											size="icon"
											@click="removeProvision(index)"
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
												class="top-2 right-2 absolute bg-gray-500 hover:bg-red-700 p-1 rounded-full text-white transition-all duration-200"
												@click="onDeleteImage"
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
						<CardDescription v-if="!readonly">Выберите категорию модификатора.</CardDescription>
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
			class="md:hidden flex justify-center items-center gap-2"
		>
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

		<!-- Dialogs -->
		<AdminIngredientsSelectDialog
			v-if="!readonly"
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<AdminSelectProvisionDialog
			:open="openProvisionsDialog"
			@close="openProvisionsDialog = false"
			@select="addProvision"
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
