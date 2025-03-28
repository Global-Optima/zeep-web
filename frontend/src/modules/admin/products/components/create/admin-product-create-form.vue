<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import type { CreateProductDTO, ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { toTypedSchema } from '@vee-validate/zod'
import { Camera, ChevronLeft, Video, X } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import {defineAsyncComponent, ref, useTemplateRef} from 'vue'
import * as z from 'zod'

const AdminSelectProductCategory = defineAsyncComponent(() =>
  import('@/modules/admin/product-categories/components/admin-select-product-category.vue'))

const {isSubmitting} = defineProps<{isSubmitting: boolean}>()

// Define emits
const emits = defineEmits<{
  onSubmit: [dto: CreateProductDTO];
  onCancel: [];
}>();

// Zod schema for validation
const createProductSchema = toTypedSchema(
  z.object({
    name: z.string()
      .min(2, 'Название должно содержать не менее 2 символов')
      .max(100, 'Название не может превышать 100 символов'),
    description: z.string()
      .max(500, 'Описание не может превышать 500 символов')
	  .optional()
	  .default(''),
    categoryId: z.coerce.number()
      .min(1, 'Выберите категорию из списка')
      .default(0),
    image: z.instanceof(File).optional().refine((file) => {
      if (!file) return true; // Optional field
      return ['image/jpeg', 'image/png'].includes(file.type);
    }, 'Поддерживаются только форматы JPEG и PNG').refine((file) => {
      if (!file) return true;
      return file.size <= 5 * 1024 * 1024; // Max 5MB
    }, 'Максимальный размер файла: 5MB'),
    video: z.instanceof(File).optional().refine((file) => {
      if (!file) return true; // Optional field
      return ['video/mp4'].includes(file.type);
    }, 'Поддерживается только формат MP4').refine((file) => {
      if (!file) return true;
      return file.size <= 20 * 1024 * 1024; // Max 20MB
    }, 'Максимальный размер файла: 20MB'),
  })
);

// Setup form with vee-validate
const { handleSubmit, setFieldValue } = useForm({
  validationSchema: createProductSchema,
});

// Submission and Cancel handlers
const onSubmit = handleSubmit((values) => {
  emits('onSubmit', values);
});
function onCancel() {
  emits('onCancel');
}

// Category selection logic
const openCategoryDialog = ref(false);
const selectedCategory = ref<ProductCategoryDTO | null>(null);
function selectCategory(category: ProductCategoryDTO) {
  selectedCategory.value = category;
  openCategoryDialog.value = false;
  setFieldValue('categoryId', category.id);
}

// File previews
const previewImage = ref<string | null>(null);
const previewVideo = ref<string | null>(null);

// File input refs
const imageInputRef = useTemplateRef("imageInputRef");
const videoInputRef = useTemplateRef("videoInputRef");

// Handle file uploads
function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    const file = target.files[0];
    setFieldValue('image', file); // Update form field
    previewImage.value = URL.createObjectURL(file); // Create preview
  }
}

function handleVideoUpload(event: Event) {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    const file = target.files[0];
    setFieldValue('video', file); // Update form field
    previewVideo.value = URL.createObjectURL(file); // Create preview
  }
}

// Trigger file input dialogs
function triggerImageInput() {
  imageInputRef.value?.click();
}

function triggerVideoInput() {
  videoInputRef.value?.click();
}
</script>

<template>
	<form
		@submit.prevent="onSubmit"
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
				Создание продукта
			</h1>
			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					:disabled="isSubmitting"
					@click="onCancel"
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

							<!-- Description -->
							<FormField
								v-slot="{ componentField }"
								name="description"
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

				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
						<CardDescription>Категория товара</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="categoryId">
							<FormItem>
								<Button
									variant="link"
									class="mt-0 p-0 h-fit text-primary underline"
									type="button"
									@click="openCategoryDialog = true"
								>
									{{ selectedCategory?.name || 'Категория не выбрана' }}
								</Button>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>
			</div>

			<!-- Right Side: Media -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Image Card -->
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

				<!-- Video Card -->
				<Card>
					<CardHeader>
						<CardTitle>Видео</CardTitle>
						<CardDescription>
							Загрузите видео для продукта.<br />
							Поддерживаемый формат: MP4 (макс. 20MB)
						</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField name="video">
							<FormItem>
								<FormControl>
									<div class="space-y-2">
										<!-- Preview -->
										<div
											v-if="previewVideo"
											class="relative rounded-lg w-full h-48 overflow-hidden"
										>
											<video
												:src="previewVideo"
												controls
												class="w-full h-full object-cover"
											></video>
											<button
												type="button"
												class="top-2 right-2 absolute bg-gray-500 transition-all duration-200 hover:bg-red-700 p-1 rounded-full text-white"
												@click="previewVideo = null; setFieldValue('video', undefined)"
											>
												<X class="size-4" />
											</button>
										</div>

										<!-- Input -->
										<div
											v-if="!previewVideo"
											class="p-4 border-2 border-gray-300 hover:border-primary border-dashed rounded-lg text-center transition-colors cursor-pointer"
											@click="triggerVideoInput"
										>
											<input
												ref="videoInputRef"
												type="file"
												accept="video/mp4"
												style="display: none;"
												@change="handleVideoUpload"
											/>
											<p class="flex flex-col justify-center items-center text-gray-500 text-sm">
												<span class="mb-2"><Video /></span>
												Нажмите для загрузки видео<br />
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
			</div>
		</div>

		<!-- Mobile Action Buttons -->
		<div class="md:hidden flex justify-center items-center gap-2">
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
		:open="openCategoryDialog"
		@close="openCategoryDialog = false"
		@select="selectCategory"
	/>
</template>
