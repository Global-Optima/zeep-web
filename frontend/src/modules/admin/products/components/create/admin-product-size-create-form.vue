<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import AdminSelectAdditiveDialog from '@/modules/admin/additives/components/admin-select-additive-dialog.vue'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import AdminIngredientsSelectDialog from '@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ProductSizeNames } from '@/modules/kiosk/products/models/product.model'
import { ChevronDown, ChevronLeft, Trash } from 'lucide-vue-next'

interface SelectedAdditiveTypesDTO {
  additiveId: number
  isDefault: boolean
  name: string
  categoryName: string
  imageUrl: string
}

interface SelectedIngredientsTypesDTO {
  ingredientId: number
  name: string
  unit: string
  category: string
  quantity: number
}

export interface CreateProductSizeFormSchema {
  name: ProductSizeNames
  unitId: number
  basePrice: number
  size: number
  machineId: string
  additives: SelectedAdditiveTypesDTO[]
  ingredients: SelectedIngredientsTypesDTO[]
}

const emits = defineEmits<{
  onSubmit: [dto: CreateProductSizeFormSchema]
  onCancel: []
}>()

const createProductSizeSchema = toTypedSchema(
  z.object({
    name: z.nativeEnum(ProductSizeNames).describe('Выберите корректный вариант'),
    basePrice: z.number().min(0, 'Введите корректную цену'),
    machineId: z.string().min(1, 'Введите код товара из автомата').max(40, "Максимум 40 символов"),
    size: z.number().min(0, 'Введите корректный размер'),
    unitId: z.number().min(1, 'Введите корректный размер'),
  })
)

const validateAdditives = (additives: SelectedAdditiveTypesDTO[]) => {
  if (!additives.length) {
    return 'Необходимо добавить хотя бы одну добавку.'
  }
  return null
}

const ingredients = ref<SelectedIngredientsTypesDTO[]>([])
const openIngredientsDialog = ref(false)

function addIngredient(ingredient: IngredientsDTO) {
  if (!ingredients.value.some((item) => item.ingredientId === ingredient.id)) {
    ingredients.value.push({
      ingredientId: ingredient.id,
      name: ingredient.name,
      unit: ingredient.unit.name,
      category: ingredient.category.name,
      quantity: 0
    })
  }
}

const { handleSubmit, isSubmitting, setFieldValue } = useForm<CreateProductSizeFormSchema>({
  validationSchema: createProductSizeSchema,
})

const additives = ref<SelectedAdditiveTypesDTO[]>([])
const additivesError = ref<string | null>(null)
const openAdditiveDialog = ref(false)

function addAdditive(additive: AdditiveDTO) {
  if (!additives.value.some((item) => item.additiveId === additive.id)) {
    additives.value.push({
      additiveId: additive.id,
      isDefault: false,
      name: additive.name,
      categoryName: additive.category.name,
      imageUrl: additive.imageUrl,
    })
  }
}

function removeAdditive(index: number) {
  additives.value.splice(index, 1)
}

function removeIngredient(index: number) {
  ingredients.value.splice(index, 1)
}

function toggleDefault(index: number) {
  additives.value[index].isDefault = !additives.value[index].isDefault
}

function sortAdditives() {
  return [...additives.value].sort((a, b) => Number(b.isDefault) - Number(a.isDefault))
}

const onSubmit = handleSubmit((formValues) => {
  additivesError.value = validateAdditives(additives.value)
  if (additivesError.value) {
    return
  }
  const finalDTO: CreateProductSizeFormSchema = {
    ...formValues,
    additives: additives.value,
    ingredients: ingredients.value
  }
  emits('onSubmit', finalDTO)
})

const onCancel = () => {
  emits('onCancel')
}

const openUnitDialog = ref(false)
const selectedUnit = ref<UnitDTO | null>(null)

function selectUnit(unit: UnitDTO) {
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}
</script>

<template>
	<div class="mx-auto w-full max-w-4xl">
		<!-- Header -->
		<div class="flex justify-between items-center gap-4 w-full">
			<div class="flex items-center gap-4">
				<Button
					variant="outline"
					size="icon"
					@click="onCancel"
				>
					<ChevronLeft class="w-5 h-5" />
					<span class="sr-only">Назад</span>
				</Button>
				<h1
					class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0"
				>
					Создать вариант
				</h1>
			</div>
			<div class="flex items-center gap-2">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					:disabled="isSubmitting"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="mt-6">
			<!-- Variant Details -->
			<Card>
				<CardHeader>
					<CardTitle>Детали варианта</CardTitle>
					<CardDescription>Укажите название, размер и цену варианта.</CardDescription>
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
									<!-- This Select uses our VariantName enum -->
									<Select v-bind="componentField">
										<SelectTrigger>
											<SelectValue placeholder="Выберите название" />
										</SelectTrigger>
										<SelectContent>
											<!-- Iterate over VariantName enum -->
											<SelectItem
												v-for="(value, key) in ProductSizeNames"
												:key="key"
												:value="value"
											>
												{{ value }}
											</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField name="unitId">
							<FormItem>
								<FormLabel>Единица измерения</FormLabel>
								<div
									@click="openUnitDialog = true"
									class="flex justify-between items-center gap-4 px-4 py-2 border rounded-md text-sm"
								>
									{{ selectedUnit?.name || 'Не выбрана' }}

									<ChevronDown class="w-5 h-5 text-gray-500" />
								</div>
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
										type="number"
										v-bind="componentField"
										placeholder="Введите размер"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<!-- Price -->
						<FormField
							name="basePrice"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Начальная цена</FormLabel>
								<FormControl>
									<Input
										type="number"
										v-bind="componentField"
										placeholder="Введите цену"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							name="machineId"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Код товара из автомата</FormLabel>
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

			<!-- Additives -->
			<Card class="mt-4">
				<CardHeader>
					<div class="flex justify-between items-start">
						<div>
							<CardTitle>Добавки</CardTitle>
							<CardDescription class="mt-2"> Выберите добавки для варианта. </CardDescription>
						</div>
						<Button
							variant="outline"
							@click="openAdditiveDialog = true"
						>
							Добавить
						</Button>
					</div>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead></TableHead>
								<TableHead>Название</TableHead>
								<TableHead>Категория</TableHead>
								<TableHead>По умолчанию</TableHead>
								<TableHead></TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="(additive, index) in sortAdditives()"
								:key="additive.additiveId"
							>
								<TableCell>
									<LazyImage
										:src="additive.imageUrl"
										alt="Изображение добавки"
										class="rounded-md size-16 object-contain"
									/>
								</TableCell>
								<TableCell>{{ additive.name }}</TableCell>
								<TableCell>{{ additive.categoryName }}</TableCell>
								<TableCell class="text-center">
									<Input
										type="checkbox"
										class="shadow-none h-5"
										:checked="additive.isDefault"
										@change="toggleDefault(index)"
									/>
								</TableCell>
								<TableCell class="text-center">
									<Button
										variant="ghost"
										size="icon"
										@click="removeAdditive(index)"
									>
										<Trash class="w-6 h-6 text-red-500" />
									</Button>
								</TableCell>
							</TableRow>
						</TableBody>
					</Table>
					<div
						v-if="additivesError"
						class="mt-2 text-red-500 text-sm"
					>
						{{ additivesError }}
					</div>
				</CardContent>
			</Card>

			<Card class="mt-4">
				<CardHeader>
					<div class="flex justify-between items-start">
						<div>
							<CardTitle>Техническая карта</CardTitle>
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
								v-for="(ingredient, index) in ingredients"
								:key="ingredient.ingredientId"
							>
								<TableCell>{{ ingredient.name }}</TableCell>
								<TableCell>{{ ingredient.category }}</TableCell>

								<TableCell class="flex items-center gap-4">
									<Input
										type="number"
										v-model.number="ingredient.quantity"
										:min="0"
										class="w-24"
										placeholder="Введите количество"
									/>
									{{ ingredient.unit.toLowerCase() }}
								</TableCell>
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
					<div
						v-if="additivesError"
						class="mt-2 text-red-500 text-sm"
					>
						{{ additivesError }}
					</div>
				</CardContent>
			</Card>
		</div>

		<AdminIngredientsSelectDialog
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<!-- Additive Dialog -->
		<AdminSelectAdditiveDialog
			:open="openAdditiveDialog"
			@close="openAdditiveDialog = false"
			@select="addAdditive"
		/>

		<AdminSelectUnit
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
	</div>
</template>
