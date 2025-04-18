<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import type { ProductCategoryDTO, ProductDetailsDTO, UpdateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { toTypedSchema } from '@vee-validate/zod'
import { Camera, ChevronLeft, Video, X } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref, useTemplateRef } from 'vue'
import * as z from 'zod'

// Lazy-load the dialog component
const AdminSelectProductCategory = defineAsyncComponent(() =>
  import('@/modules/admin/product-categories/components/admin-select-product-category.vue')
);
const AdminProductsTranslationsDialog = defineAsyncComponent(() =>
  import('@/modules/admin/products/components/admin-products-translations-dialog.vue')
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

const updateProductSchema = toTypedSchema(
  z.object({
    name: z.string()
      .min(2, 'Название должно содержать не менее 2 символов')
      .max(100, 'Название не может превышать 100 символов'),
    description: z.string()
      .max(500, 'Описание не может превышать 500 символов')
	  .optional(),
    categoryId: z.coerce.number().min(1, 'Выберите категорию из списка'),
    image: z.instanceof(File).optional().refine((file) => {
      if (!file) return true;
      return ['image/jpeg', 'image/png'].includes(file.type);
    }, 'Поддерживаются только форматы JPEG и PNG').refine((file) => {
      if (!file) return true;
      return file.size <= 5 * 1024 * 1024;
    }, 'Максимальный размер файла: 5MB'),
    video: z.instanceof(File).optional().refine((file) => {
      if (!file) return true;
      return ['video/mp4'].includes(file.type);
    }, 'Поддерживается только формат MP4').refine((file) => {
      if (!file) return true;
      return file.size <= 20 * 1024 * 1024;
    }, 'Максимальный размер файла: 20MB'),
  })
);

const { handleSubmit, setFieldValue, resetField } = useForm({
  validationSchema: updateProductSchema,
  initialValues: {
    name: productDetails.name,
    description: productDetails.description,
    categoryId: productDetails.category.id,
  },
});

const openCategoryDialog = ref(false);
const selectedCategory = ref<ProductCategoryDTO | null>(productDetails.category);

function selectCategory(category: ProductCategoryDTO) {
  if (!readonly) {
    selectedCategory.value = category;
    openCategoryDialog.value = false;
    setFieldValue('categoryId', category.id);
  }
}

const previewImage = ref<string | null>(productDetails.imageUrl || null);
const previewVideo = ref<string | null>(productDetails.videoUrl || null);
const deleteImage = ref<boolean>(false)
const deleteVideo = ref<boolean>(false)

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.length) {
    const file = target.files[0];
    previewImage.value = URL.createObjectURL(file);
    setFieldValue('image', file);
    deleteImage.value = true
  }
};

const handleVideoUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.length) {
    const file = target.files[0];
    previewVideo.value = URL.createObjectURL(file);
    setFieldValue('video', file);
    deleteVideo.value = true
  }
};

const onSubmit = handleSubmit((values) => {
  const dto: UpdateProductDTO = {
    name: values.name,
    description: values.description,
    categoryId: values.categoryId,
    image: values.image,
    video: values.video,
    deleteImage: deleteImage.value,
    deleteVideo: deleteVideo.value,
  };

  emits('onSubmit', dto);
});

const resetFormValues = () => {
  resetField('image')
  resetField('video')
  deleteImage.value = false
  deleteVideo.value = false;
};

defineExpose({ resetFormValues });

const onDeleteImage = () => {
  previewImage.value = null
  setFieldValue('image', undefined)
  deleteImage.value = true
}

const onDeleteVideo = () => {
  previewVideo.value = null
  setFieldValue('video', undefined)
  deleteVideo.value = true
}

const onCancel = () => emits('onCancel');

const imageInputRef = useTemplateRef("imageInputRef");
const videoInputRef = useTemplateRef("videoInputRef");
const triggerImageInput = () => imageInputRef.value?.click();
const triggerVideoInput = () => videoInputRef.value?.click();
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
				<AdminProductsTranslationsDialog :productId="productDetails.id" />

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
						<CardDescription>Категория продукта</CardDescription>
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
												class="top-2 right-2 absolute bg-gray-500 hover:bg-red-700 p-1 rounded-full text-white transition-all duration-200"
												@click="onDeleteVideo"
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
