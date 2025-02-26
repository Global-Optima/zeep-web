<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import type { ProductCategoryDTO, ProductDetailsDTO, UpdateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronLeft } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

// Lazy-load the dialog component
const AdminSelectProductCategory = defineAsyncComponent(() =>
  import('@/modules/admin/product-categories/components/admin-select-product-category.vue')
);

const {productDetails, readonly, isSubmitting} = defineProps<{
  productDetails: ProductDetailsDTO;
  readonly?: boolean;
  isSubmitting: boolean
}>();

const emits = defineEmits<{
  onSubmit: [dto: UpdateProductDTO];
  onCancel: [];
}>();

const createProductSchema = toTypedSchema(
  z.object({
    name: z.string()
      .min(2, 'Название должно содержать не менее 2 символов')
      .max(100, 'Название не может превышать 100 символов'),
    description: z.string()
      .max(500, 'Описание не может превышать 500 символов'),
    categoryId: z.coerce.number()
      .min(1, 'Выберите категорию из списка'),
    imageUrl: z.string().min(1, 'Вставьте картинку добавки'),
  })
);

const { handleSubmit, isFieldDirty, setFieldValue, resetForm } = useForm<UpdateProductDTO>({
  validationSchema: createProductSchema,
  initialValues: {
    name: productDetails.name,
    description: productDetails.description,
    categoryId: productDetails.category.id,
    imageUrl: productDetails.imageUrl,
  },
});

const onSubmit = handleSubmit((values) => {
  emits('onSubmit', values);
});

function onCancel() {
  resetForm();
  emits('onCancel');
}

const openCategoryDialog = ref(false);
const selectedCategory = ref<ProductCategoryDTO | null>(productDetails.category);

function selectCategory(category: ProductCategoryDTO) {
  if (!readonly) {
    selectedCategory.value = category;
    openCategoryDialog.value = false;
    setFieldValue('categoryId', category.id);
  }
}
</script>

<template>
	<form
		@submit="onSubmit"
		class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl"
	>
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				type="button"
				@click="onCancel"
				:disabled="isSubmitting"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				{{ productDetails.name }}
			</h1>
			<div
				class="hidden md:flex items-center gap-2 md:ml-auto"
				v-if="!readonly"
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
					:disabled="isSubmitting"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<!-- Left Side: Product Details -->
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Основная информация</CardTitle>
						<CardDescription>Заполните основные сведения о продукте.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<!-- Name -->
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
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Description -->
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
											:readonly="readonly"
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
						<CardDescription>Категория товара</CardDescription>
					</CardHeader>
					<CardContent>
						<div>
							<template v-if="!readonly">
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-primary underline"
									type="button"
									@click="openCategoryDialog = true"
									:disabled="readonly"
								>
									{{ selectedCategory?.name || 'Категория не выбрана' }}
								</Button>
							</template>
							<template v-else>
								<span class="text-muted-foreground">
									{{ selectedCategory?.name || 'Категория не выбрана' }}
								</span>
							</template>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Right Side: Media -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Image Preview Card -->
				<Card>
					<CardHeader>
						<CardTitle>Изображение</CardTitle>
						<CardDescription> Предварительный просмотр изображения продукта. </CardDescription>
					</CardHeader>
					<CardContent>
						<div class="space-y-2">
							<div
								v-if="productDetails.imageUrl"
								class="relative border rounded-lg w-full h-48 overflow-hidden"
							>
								<LazyImage
									:src="productDetails.imageUrl"
									alt="Product Image"
									class="rounded-lg w-full h-full object-contain"
								/>
							</div>
							<div
								v-else
								class="p-4 border-2 border-gray-300 border-dashed rounded-lg text-center"
							>
								<p class="flex flex-col justify-center items-center text-gray-500 text-sm">
									<span class="mb-2">📷</span>
									Изображение отсутствует
								</p>
							</div>
						</div>
					</CardContent>
				</Card>

				<!-- Video Preview Card -->
				<Card>
					<CardHeader>
						<CardTitle>Видео</CardTitle>
						<CardDescription> Предварительный просмотр видео продукта. </CardDescription>
					</CardHeader>
					<CardContent>
						<div class="space-y-2">
							<div
								v-if="productDetails.videoUrl"
								class="relative rounded-lg w-full h-48 overflow-hidden"
							>
								<video
									:src="productDetails.videoUrl"
									controls
									class="w-full h-full object-cover"
								></video>
							</div>
							<div
								v-else
								class="p-4 border-2 border-gray-300 border-dashed rounded-lg text-center"
							>
								<p class="flex flex-col justify-center items-center text-gray-500 text-sm">
									<span class="mb-2">🎥</span>
									Видео отсутствует
								</p>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Mobile Action Buttons -->
		<div
			class="md:hidden flex justify-center items-center gap-2"
			v-if="!readonly"
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
				:disabled="isSubmitting"
				>Сохранить</Button
			>
		</div>
	</form>

	<AdminSelectProductCategory
		v-if="!readonly"
		:open="openCategoryDialog"
		@close="openCategoryDialog = false"
		@select="selectCategory"
	/>
</template>
