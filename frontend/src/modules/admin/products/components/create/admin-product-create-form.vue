<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'

// Icons
import { ChevronLeft } from 'lucide-vue-next'

// Services / Models
import AdminSelectProductCategory from '@/modules/admin/product-categories/components/admin-select-product-category.vue'
import type { CreateProductDTO, ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { ref } from 'vue'

// 2. Define emits
const emits = defineEmits<{
  onSubmit: [dto: CreateProductDTO]
  onCancel: []
}>()

// 3. Zod schema for the CreateProductDTO
//    (You can remove videoUrl if you are not actually using it.)
const createProductSchema = toTypedSchema(
  z.object({
    name: z.string()
      .min(2, 'Название должно содержать не менее 2 символов')
      .max(100, 'Название не может превышать 100 символов'),
    description: z.string()
      .max(500, 'Описание не может превышать 500 символов'),
    imageUrl: z.string()
      .url('Введите корректную ссылку (URL)'),
    categoryId: z.coerce.number()
      .min(1, 'Выберите категорию из списка'),
  })
)

// 4. Setup form with vee-validate
const { handleSubmit, isSubmitting, isFieldDirty, setFieldValue } = useForm<CreateProductDTO>({
  validationSchema: createProductSchema,
})

// 5. Submission and Cancel handlers
const onSubmit = handleSubmit((values) => {
  emits('onSubmit', values)
})

function onCancel() {
  emits('onCancel')
}

const openCategoryDialog = ref(false)
const selectedCategory = ref<ProductCategoryDTO | null>(null)

function selectCategory(category: ProductCategoryDTO) {
  selectedCategory.value = category
  openCategoryDialog.value = false
  setFieldValue('categoryId', category.id)
}
</script>

<template>
	<form
		@submit.prevent="onSubmit"
		class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl"
	>
		<!-- ========== Header ========== -->
		<div class="flex items-center gap-4">
			<!-- Back / Cancel button -->
			<Button
				variant="outline"
				size="icon"
				type="button"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>

			<!-- Title -->
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Создание продукта
			</h1>

			<!-- Desktop action buttons -->
			<div class="md:flex items-center gap-2 hidden md:ml-auto">
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
				>
					Сохранить
				</Button>
			</div>
		</div>

		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<!-- LEFT side: Product Details (Name, Description) -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Основная информация</CardTitle>
						<CardDescription>
							Заполните основные сведения о продукте (название и описание).
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<!-- Название -->
							<FormField
								v-slot="{ componentField }"
								name="name"
								:validate-on-blur="!isFieldDirty"
							>
								<FormItem>
									<FormLabel>Название продукта</FormLabel>
									<FormControl>
										<Input
											type="text"
											placeholder="Например, Латте"
											v-bind="componentField"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Описание -->
							<FormField
								v-slot="{ componentField }"
								name="description"
								:validate-on-blur="!isFieldDirty"
							>
								<FormItem>
									<FormLabel>Описание</FormLabel>
									<FormControl>
										<Textarea
											placeholder="Краткое описание продукта"
											class="min-h-32"
											v-bind="componentField"
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
				<!-- Media Card -->
				<Card>
					<CardHeader>
						<CardTitle>Медиа</CardTitle>
						<CardDescription>
							Вставьте ссылки на изображение или видео для продукта.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-4 grid">
							<!-- Ссылка на изображение -->
							<FormField
								v-slot="{ componentField }"
								name="imageUrl"
								:validate-on-blur="!isFieldDirty"
							>
								<FormItem>
									<FormLabel>Ссылка на изображение</FormLabel>
									<FormControl>
										<Input
											type="text"
											placeholder="https://example.com/image.jpg"
											v-bind="componentField"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>

				<!-- Category Card -->
				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription>Выберите категорию товара</CardDescription>
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

		<!-- Mobile action buttons (only visible on small screens) -->
		<div class="flex justify-center items-center gap-2 md:hidden">
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
			>
				Сохранить
			</Button>
		</div>
	</form>

	<AdminSelectProductCategory
		:open="openCategoryDialog"
		@close="openCategoryDialog = false"
		@select="selectCategory"
	/>
</template>
