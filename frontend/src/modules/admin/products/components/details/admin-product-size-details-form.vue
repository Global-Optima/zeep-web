<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

// UI Components
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import Checkbox from '@/core/components/ui/checkbox/Checkbox.vue'
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
import { TechnicalMapEntity, useCopyTechnicalMap } from '@/core/hooks/use-copy-technical-map.hooks'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ProductSizeNames, type ProductSizeDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronDown, ChevronLeft, EllipsisVertical, Trash } from 'lucide-vue-next'
import type { ProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import { getRouteName } from "@/core/config/routes.config"

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

export interface UpdateProductSizeFormSchema {
  name: ProductSizeNames
  unitId: number
  basePrice: number
  size: number
  machineId: string
  additives: SelectedAdditiveTypesDTO[]
  ingredients: SelectedIngredientsTypesDTO[]
  provisions: SelectedProvisionsTypesDTO[]
}

const { productSize, readonly = false } = defineProps<{
  productSize: ProductSizeDetailsDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateProductSizeFormSchema]
  onCancel: []
}>()

const { setTechnicalMapReference, fetchTechnicalMap } = useCopyTechnicalMap()

const updateProductSizeSchema = toTypedSchema(
  z.object({
    name: z.nativeEnum(ProductSizeNames).describe('Выберите корректный вариант'),
    basePrice: z.number().min(0, 'Введите корректную цену'),
    size: z.number().min(0, 'Введите корректный размер'),
    unitId: z.number().min(1, 'Введите корректную единицу измерения'),
    machineId: z.string().min(1, 'Введите код продукта из автомата').max(40, "Максимум 40 символов"),
  })
)

const { handleSubmit, isSubmitting, setFieldValue } = useForm({
  validationSchema: updateProductSizeSchema,
  initialValues: {
    name: productSize.name as ProductSizeNames,
    unitId: productSize.unit.id,
    basePrice: productSize.basePrice,
    size: productSize.size,
    machineId: productSize.machineId
  }
})

const provisions = ref<SelectedProvisionsTypesDTO[]>(productSize.provisions.map(p => ({
  provisionId: p.provision.id,
  name: p.provision.name,
  absoluteVolume: p.provision.absoluteVolume,
  unit: p.provision.unit.name,
  volume: p.volume,
})))
const openProvisionsDialog = ref(false)

const additives = ref<SelectedAdditiveTypesDTO[]>(productSize.additives.map(a => ({
  additiveId: a.id,
  isDefault: a.isDefault,
  isHidden: a.isHidden,
  name: a.name,
  categoryName: a.category.name,
  size: a.size,
  unitName: a.unit.name,
  imageUrl: a.imageUrl
})))

const additivesError = ref<string | null>(null)
const openAdditiveDialog = ref(false)

const ingredients = ref<SelectedIngredientsTypesDTO[]>(productSize.ingredients.map(i => ({
  ingredientId: i.ingredient.id,
  name: i.ingredient.name,
  unit: i.ingredient.unit.name,
  category: i.ingredient.category.name,
  quantity: i.quantity,
})))

const openIngredientsDialog = ref(false)

function addProvision(provision: ProvisionDTO) {
  if (readonly) return
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

function addAdditive(additive: AdditiveDTO) {
  if (readonly) return
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

function addIngredient(ingredient: IngredientsDTO) {
  if (readonly) return
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

function removeAdditive(index: number) {
  if (readonly) return
  additives.value.splice(index, 1)
}

function removeIngredient(index: number) {
  if (readonly) return
  ingredients.value.splice(index, 1)
}

function removeProvision(index: number) {
  if (readonly) return
  provisions.value.splice(index, 1)
}

const onSubmit = handleSubmit((formValues) => {
  if (readonly) return

  if (additivesError.value) {
    return
  }

  if (ingredients.value.some(i => i.quantity <= 0)) {
    return toast({ description: "Укажите количество в технологической карте" })
  }

  if (provisions.value.some(i => i.volume <= 0)) {
    return toast({ description: "Укажите обьем в заготовке" })
  }

  const finalDTO: UpdateProductSizeFormSchema = {
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
const selectedUnit = ref<UnitDTO | null>(productSize.unit)

function selectUnit(unit: UnitDTO) {
  if (readonly) return
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

const onCopyTechMapClick = () => {
  try {
    setTechnicalMapReference(TechnicalMapEntity.PRODUCT_SIZE, productSize.id)
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
					Детали варианта
				</h1>
			</div>
			<div
				v-if="!readonly"
				class="flex items-center gap-2"
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
					<CardDescription v-if="!readonly">
						Укажите название, размер и цену варианта.
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
									<Select
										v-bind="componentField"
										:disabled="readonly"
									>
										<SelectTrigger>
											<SelectValue placeholder="Выберите название" />
										</SelectTrigger>
										<SelectContent>
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
							<FormItem class="flex flex-col gap-1">
								<FormLabel>Единица измерения</FormLabel>
								<FormControl>
									<div
										@click="!readonly && (openUnitDialog = true)"
										class="flex justify-between items-center gap-4 px-4 py-2 border rounded-md text-sm"
										:class="{ 'cursor-pointer': !readonly }"
									>
										{{ selectedUnit?.name || 'Не выбрана' }}
										<ChevronDown
											v-if="!readonly"
											class="w-5 h-5 text-gray-500"
										/>
									</div>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<!-- Size -->
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
										:readonly="readonly"
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
										:readonly="readonly"
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
							<CardDescription
								v-if="!readonly"
								class="mt-2"
							>
								Выберите модификатора для варианта.
							</CardDescription>
						</div>
						<Button
							v-if="!readonly"
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
								<TableHead class="!pr-2 text-center">В составе</TableHead>
								<TableHead class="!pr-2 text-center">Не показывать</TableHead>
								<TableHead v-if="!readonly"></TableHead>
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
								<TableCell>
									<RouterLink
										:to="{name: getRouteName('ADMIN_ADDITIVE_DETAILS'), params: {id: additive.additiveId}}"
										target="_blank"
										class="hover:text-primary underline transition-colors duration-300 underline-offset-4"
									>
										{{ additive.name }}
									</RouterLink>
								</TableCell>
								<TableCell>{{ additive.categoryName }}</TableCell>
								<TableCell>{{ additive.size }} {{ additive.unitName }}</TableCell>
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

								<TableCell
									v-if="!readonly"
									class="text-center"
								>
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
								<TableHead>Количество</TableHead>
								<TableHead>Размер</TableHead>
								<TableHead v-if="!readonly"></TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="(ingredient, index) in ingredients"
								:key="ingredient.ingredientId"
							>
								<TableCell>
									<RouterLink
										:to="{name: getRouteName('ADMIN_INGREDIENTS_DETAILS'), params: {id: ingredient.ingredientId}}"
										target="_blank"
										class="hover:text-primary underline transition-colors duration-300 underline-offset-4"
									>
										{{ ingredient.name }}
									</RouterLink>
								</TableCell>

								<TableCell class="flex items-center gap-2 w-24">
									<Input
										type="number"
										v-model.number="ingredient.quantity"
										:min="0"
										placeholder="Введите количество"
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
								<TableHead></TableHead>
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
										class="hover:text-primary underline transition-colors duration-300 underline-offset-4"
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

		<!-- Dialogs -->
		<AdminSelectAdditiveDialog
			v-if="!readonly"
			:open="openAdditiveDialog"
			@close="openAdditiveDialog = false"
			@select="addAdditive"
		/>

		<AdminSelectProvisionDialog
			:open="openProvisionsDialog"
			@close="openProvisionsDialog = false"
			@select="addProvision"
		/>

		<AdminIngredientsSelectDialog
			v-if="!readonly"
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>

		<AdminSelectUnit
			v-if="!readonly"
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>
	</div>
</template>
