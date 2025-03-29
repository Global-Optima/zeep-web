<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'
import { ref } from 'vue';

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import type { ProductCategoryDTO, UpdateProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft } from 'lucide-vue-next'

const { productCategory, readonly = false } = defineProps<{
  productCategory: ProductCategoryDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateProductCategoryDTO]
  onCancel: []
}>()

// Validation Schema
const createCategorySchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название категории'),
    description: z.string().min(1, 'Введите описание категории'),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm<UpdateProductCategoryDTO>({
  validationSchema: createCategorySchema,
  initialValues: productCategory
})

const machineCategoryOptions = ref([
  { label: "Чай", value: "TEA" },
  { label: "Кофе", value: "COFFEE" },
  { label: "Мороженое", value: "ICE_CREAM" }
]);

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (readonly) return
  emits('onSubmit', formValues)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
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
				{{ productCategory.name }}
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
		<Card>
			<CardHeader>
				<CardTitle>Детали категории</CardTitle>
				<CardDescription v-if="!readonly">Заполните информацию о категории.</CardDescription>
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
									placeholder="Введите название категории"
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
								<Input
									id="description"
									type="text"
									v-bind="componentField"
									placeholder="Введите описание категории"
									:readonly="readonly"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Machine Category -->
					<FormField
						name="machineCategory"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Категория Машины</FormLabel>
							<FormControl>
								<Select v-bind="componentField">
									<SelectTrigger id="machine_category">
										<SelectValue placeholder="Выберите категорию машины" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem
											v-for="option in machineCategoryOptions"
											:key="option.value"
											:value="option.value"
										>
											{{ option.label }}
										</SelectItem>
									</SelectContent>
								</Select>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Mobile Footer -->
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
	</div>
</template>
