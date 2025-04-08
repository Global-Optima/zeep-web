<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref, toRefs } from 'vue'
import * as z from 'zod'

// UI Components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import Checkbox from "@/core/components/ui/checkbox/Checkbox.vue"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/core/components/ui/dropdown-menu'
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
import { toast } from "@/core/components/ui/toast"
import { useCopyTechnicalMap } from '@/core/hooks/use-copy-technical-map.hooks'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { ProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ProductSizeNames, type ProductSizeDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronDown, ChevronLeft, EllipsisVertical, Trash } from 'lucide-vue-next'
import { watch } from 'vue'

const AdminSelectAdditiveDialog = defineAsyncComponent(() =>
  import('@/modules/admin/additives/components/admin-select-additive-dialog.vue'))
const AdminIngredientsSelectDialog = defineAsyncComponent(() =>
  import('@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue'))
const AdminSelectUnit = defineAsyncComponent(() =>
  import('@/modules/admin/units/components/admin-select-unit.vue'))
const AdminSelectProvisionDialog = defineAsyncComponent(() =>
  import("@/modules/admin/provisions/components/admin-select-provision-dialog.vue"))

interface SelectedAdditiveTypesDTO {
  additiveId: number
  isDefault: boolean
  isHidden: boolean
  name: string
  categoryName: string
  size: number
  unitName: string
  imageUrl: string
}

interface SelectedIngredientsTypesDTO {
  ingredientId: number
  name: string
  unit: string
  category: string
  quantity: number
}

interface SelectedProvisionsTypesDTO {
  provisionId: number
  name: string
  absoluteVolume: number
  unit: string
  volume: number
}

export interface CreateProductSizeFormSchema {
  name: ProductSizeNames
  unitId: number
  basePrice: number
  size: number
  machineId: string
  additives: SelectedAdditiveTypesDTO[]
  ingredients: SelectedIngredientsTypesDTO[]
  provisions: SelectedProvisionsTypesDTO[]
}

const props = defineProps<{initialProductSize?: ProductSizeDetailsDTO }>()
const { initialProductSize } = toRefs(props)


const emits = defineEmits<{
  onSubmit: [dto: CreateProductSizeFormSchema]
  onCancel: []
}>()

const createProductSizeSchema = toTypedSchema(
  z.object({
    name: z.nativeEnum(ProductSizeNames).describe('Выберите корректный вариант'),
    basePrice: z.number().min(0, 'Введите корректную цену'),
    machineId: z.string().min(1, 'Введите код продукта из автомата').max(40, "Максимум 40 символов"),
    size: z.number().min(0, 'Введите корректный размер'),
    unitId: z.number().min(1, 'Введите корректный размер'),
  })
)

const ingredients = ref<SelectedIngredientsTypesDTO[]>(initialProductSize.value?.ingredients?.map(i => ({
  ingredientId: i.ingredient.id,
  name: i.ingredient.name,
  unit: i.ingredient.unit.name,
  category: i.ingredient.category.name,
  quantity: i.quantity
})) || [])

const openIngredientsDialog = ref(false)

const provisions = ref<SelectedProvisionsTypesDTO[]>([])
const openProvisionsDialog = ref(false)

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

const { handleSubmit, isSubmitting, setFieldValue, resetForm } = useForm({
  validationSchema: createProductSizeSchema,
  initialValues: {
	name: initialProductSize.value?.name,
	basePrice: initialProductSize.value?.basePrice,
	size: initialProductSize.value?.size,
	unitId: initialProductSize.value?.unit.id,
	machineId: initialProductSize.value?.machineId,
  },
})

const { fetchTechnicalMap } = useCopyTechnicalMap()

const additives = ref<SelectedAdditiveTypesDTO[]>(initialProductSize.value?.additives?.map(a => ({
  additiveId: a.id,
  isDefault: a.isDefault,
  isHidden: a.isHidden,
  name: a.name,
  categoryName: a.category.name,
  size: a.size,
  unitName: a.unit.name,
  imageUrl: a.imageUrl,
})) || [])
const additivesError = ref<string | null>(null)
const openAdditiveDialog = ref(false)

function addAdditive(additive: AdditiveDTO) {
  if (!additives.value.some((item) => item.additiveId === additive.id)) {
    additives.value.push({
      additiveId: additive.id,
      isDefault: false,
      isHidden: false,
      name: additive.name,
      categoryName: additive.category.name,
      size: additive.size,
      unitName: additive.unit.name,
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

function removeProvision(index: number) {
  provisions.value.splice(index, 1)
}

const onSubmit = handleSubmit((formValues) => {
  if (additivesError.value) {
    return
  }

  if (ingredients.value.some(i => i.quantity <= 0)) {
    return toast({ description: "Укажите количество в технологической карте" })
  }

  if (provisions.value.some(i => i.volume <= 0)) {
    return toast({ description: "Укажите обьем в заготовке" })
  }

  const finalDTO: CreateProductSizeFormSchema = {
    ...formValues,
    additives: additives.value,
    ingredients: ingredients.value,
    provisions: provisions.value
  }

  emits('onSubmit', finalDTO)
})

const onCancel = () => {
  emits('onCancel')
}

const openUnitDialog = ref(false)
const selectedUnit = ref<UnitDTO | null>(initialProductSize.value?.unit || null)

function selectUnit(unit: UnitDTO) {
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}

function onAdditiveDefaultClick(index: number, value: boolean) {
  additives.value[index].isDefault = value
  if (!value) {
    additives.value[index].isHidden = false
  }
}

const onPasteTechMapClick = async () => {
  try {
    const techMap = await fetchTechnicalMap()
    if (!techMap) {
      toast({ description: "Технологическая карта не найдена" })
      return
    }

    ingredients.value = techMap.map(t => ({
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

watch(initialProductSize, (newVal) => {
  if (newVal) {
    // Reset the form with new initial values
    resetForm({
      values: {
        name: newVal.name,
        basePrice: newVal.basePrice,
        size: newVal.size,
        unitId: newVal.unit.id,
        machineId: newVal.machineId,
      }
    })
    // Update reactive states
    ingredients.value = newVal.ingredients.map(i => ({
      ingredientId: i.ingredient.id,
      name: i.ingredient.name,
      unit: i.ingredient.unit.name,
      category: i.ingredient.category.name,
      quantity: i.quantity
    }))

    additives.value = newVal.additives.map(a => ({
      additiveId: a.id,
      isDefault: a.isDefault,
      isHidden: a.isHidden,
      name: a.name,
      categoryName: a.category.name,
      size: a.size,
      unitName: a.unit.name,
      imageUrl: a.imageUrl,
    }))

    provisions.value = newVal.provisions.map(p => ({
      provisionId: p.provision.id,
      name: p.provision.name,
      absoluteVolume: p.provision.absoluteVolume,
      unit: p.provision.unit.name,
      volume: p.volume,
    }))

    selectedUnit.value = newVal.unit
  }
}, { immediate: true })
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
								<FormLabel>Код продукта из автомата</FormLabel>
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
							<CardTitle>Модификаторы</CardTitle>
							<CardDescription class="mt-2"> Выберите модификатора для варианта. </CardDescription>
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
								<TableHead>Размер</TableHead>
								<TableHead class="text-center">В составе</TableHead>
								<TableHead class="text-center">Не показывать</TableHead>
								<TableHead></TableHead>
							</TableRow>
						</TableHeader>

						<TableBody>
							<TableRow
								v-for="(additive, index) in additives"
								:key="additive.additiveId"
							>
								<TableCell>
									<LazyImage
										:src="additive.imageUrl"
										alt="Изображение модификатора"
										class="rounded-md size-16 object-contain"
									/>
								</TableCell>

								<TableCell>{{ additive.name }}</TableCell>
								<TableCell>{{ additive.categoryName }}</TableCell>
								<TableCell>{{ additive.size }} {{ additive.unitName }}</TableCell>

								<!-- По умолчанию -->
								<TableCell class="!pr-2 text-center">
									<Checkbox
										type="checkbox"
										class="data-[state=checked]:bg-slate-500 border-slate-400 size-6 data-[state=checked]:text-white"
										:checked="additive.isDefault"
										@update:checked="value => onAdditiveDefaultClick(index, value)"
									/>
								</TableCell>

								<TableCell class="!pr-2 text-center">
									<Checkbox
										type="checkbox"
										:disabled="!additive.isDefault"
										class="data-[state=checked]:bg-slate-500 border-slate-400 size-6 data-[state=checked]:text-white"
										:checked="additive.isHidden ?? false"
										@update:checked="v => additive.isHidden = v"
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

						<div class="flex items-center gap-2">
							<DropdownMenu>
								<DropdownMenuTrigger class="p-2 border rounded-md">
									<EllipsisVertical class="size-4" />
								</DropdownMenuTrigger>
								<DropdownMenuContent>
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
				</CardContent>
			</Card>

			<Card class="mt-4">
				<CardHeader>
					<div class="flex justify-between items-start">
						<div>
							<CardTitle>Заготовки</CardTitle>
							<CardDescription class="mt-2"> Выберите заготовки и их обьем </CardDescription>
						</div>
						<Button
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
								<TableHead></TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="(provision, index) in provisions"
								:key="provision.provisionId"
							>
								<TableCell>{{ provision.name }}</TableCell>
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
								<TableCell class="text-center">
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

		<AdminIngredientsSelectDialog
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<AdminSelectProvisionDialog
			:open="openProvisionsDialog"
			@close="openProvisionsDialog = false"
			@select="addProvision"
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
