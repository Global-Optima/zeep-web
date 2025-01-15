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
import AdminSelectAdditiveCategory from '@/modules/admin/additive-categories/components/admin-select-additive-category.vue'
import type { AdditiveCategoryDTO, AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { ChevronLeft } from 'lucide-vue-next'

// Props
const props = defineProps<{
  additive: AdditiveDTO
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateAdditiveFormSchema]
  onCancel: []
}>()

// DTO
export interface UpdateAdditiveFormSchema {
  name?: string
  description?: string
  price?: number
  imageUrl?: string
  size?: string
  additiveCategoryId?: number
}

// Reactive State for Category Selection
const selectedCategory = ref(props.additive.category)
const openCategoryDialog = ref(false)

// Validation Schema
const updateAdditiveSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название добавки'),
    description: z.string().min(1, 'Введите описание добавки'),
    imageUrl: z.string().min(1, 'Введите ссылку на изображение добавки'),
    price: z.number().min(0, 'Введите корректную цену'),
    size: z.string().min(0, 'Введите корректный размер'),
    additiveCategoryId: z.coerce.number().min(1, 'Выберите категорию добавки'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm<UpdateAdditiveFormSchema>({
  validationSchema: updateAdditiveSchema,
  initialValues: {
    name: props.additive.name || '',
    description: props.additive.description || '',
    price: props.additive.basePrice || 0,
    imageUrl: props.additive.imageUrl || '',
    size: props.additive.size || '',
    additiveCategoryId: props.additive.category.id || undefined,
  },
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  emits('onSubmit', { ...formValues, additiveCategoryId: selectedCategory.value.id })
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}

function selectCategory(category: AdditiveCategoryDTO) {
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue("additiveCategoryId", category.id)
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
			<!-- Additive Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали добавки</CardTitle>
						<CardDescription>Обновите название, описание и цену добавки.</CardDescription>
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
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Price and Size -->
							<div class="flex gap-4">
								<FormField
									name="price"
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
												placeholder="Например: мл, гр"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Media and Category Blocks -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Media Block -->
				<Card>
					<CardHeader>
						<CardTitle>Медиа</CardTitle>
						<CardDescription>Загрузите изображение добавки.</CardDescription>
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
						<CardDescription>Выберите категорию топпинга</CardDescription>
					</CardHeader>
					<CardContent>
						<div>
							<Button
								variant="link"
								class="mt-0 p-0 h-fit text-blue-600 underline"
								@click="openCategoryDialog = true"
							>
								{{ selectedCategory?.name || 'Категория не выбрана' }}
							</Button>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
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

		<!-- Category Dialog -->
		<AdminSelectAdditiveCategory
			:open="openCategoryDialog"
			@close="openCategoryDialog = false"
			@select="selectCategory"
		/>
	</div>
</template>
