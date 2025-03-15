<template>
	<div class="relative bg-[#F3F4F9] pb-32 text-black no-scrollbar">
		<!-- Loading State -->
		<PageLoader v-if="isFetching" />
		<!-- Error State -->
		<div
			v-else-if="isError"
			class="flex flex-col justify-center items-center p-6 w-full h-screen text-center"
		>
			<h1 class="mt-8 font-bold text-red-600 text-4xl">Ошибка</h1>
			<p class="mt-6 max-w-md text-gray-600 text-2xl">
				К сожалению, данный товар временно недоступен, попробуйте позже
			</p>

			<button
				@click="onBackClick"
				class="flex justify-center items-center bg-slate-200 mt-8 px-8 py-5 rounded-3xl h-14 text-slate- text-2xl"
			>
				Вернуться назад
			</button>
		</div>
		<!-- Product Content -->
		<div
			v-else-if="productDetails"
			class="pb-44"
		>
			<div class="bg-white shadow-gray-200 shadow-xl px-8 pb-6 rounded-b-[48px] w-full">
				<header class="flex justify-between items-center gap-6 pt-6">
					<Button
						size="icon"
						variant="ghost"
						class="size-14"
						@click="onBackClick"
					>
						<ChevronLeft
							class="size-14 text-gray-400"
							stroke-width="1.6"
						/>
					</Button>

					<template v-if="selectedSize">
						<KioskProductRecipeDialog :nutrition="selectedSize.totalNutrition" />
					</template>
				</header>
				<div class="flex flex-col justify-center items-center">
					<LazyImage
						src="https://www.nicepng.com/png/full/106-1060376_starbucks-iced-coffee-png-vector-library-pumpkin-spice.png"
						alt="Изображение товара"
						class="w-38 h-64 object-contain"
					/>
					<p class="mt-7 font-semibold text-4xl">{{ productDetails.name }}</p>
					<p class="mt-2 text-slate-600 text-xl">{{ productDetails.description }}</p>
				</div>
				<!-- Sticky Section -->
				<div class="top-0 z-10 sticky bg-white mt-8 pb-6">
					<div class="flex justify-between items-center gap-4">
						<!-- Size Selection -->
						<div class="flex items-center gap-2 overflow-x-auto no-scrollbar">
							<KioskDetailsSizes
								v-for="size in sortedSizes"
								:key="size.id"
								:size="size"
								:is-selected="selectedSize?.id === size.id"
								@click:size="onSizeSelect"
							/>
						</div>
						<!-- Add to Cart Button -->
						<div class="flex items-center gap-6">
							<p class="font-medium text-4xl">{{ formatPrice(totalPrice) }}</p>
							<button
								@click="handleAddToCart"
								class="flex items-center gap-3 bg-primary p-5 rounded-full text-primary-foreground"
							>
								<Plus class="size-8 sm:size-12" />
							</button>
						</div>
					</div>
				</div>
			</div>
			<!-- Additives Selection -->
			<div class="mt-10">
				<KioskDetailsAdditivesSection
					:categories="additiveCategories ?? []"
					:isAdditiveSelected="isAdditiveSelected"
					@toggle-additive="onAdditiveToggle"
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type {
  StoreAdditiveCategoryDTO,
  StoreAdditiveCategoryItemDTO
} from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, Plus } from 'lucide-vue-next'
import { computed, ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();

const productId = computed(() => Number(route.params.id));

const onBackClick = () => router.push({ name: getRouteName('KIOSK_HOME') });

const selectedSize = ref<StoreProductSizeDetailsDTO | null>(null);
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>({});

const { data: productDetails, isPending: isFetching, isError, error } = useQuery({
  queryKey: computed(() => ['kiosk-product-details', productId]),
  queryFn: () => storeProductsService.getStoreProduct(productId.value),
  enabled: computed(() => productId.value > 0),
});

const sortedSizes = computed(() => {
  return productDetails.value?.sizes ? [...productDetails.value.sizes].sort((a, b) => a.size - b.size) : [];
});

watchEffect(() => {
  if (!selectedSize.value && sortedSizes.value.length > 0) {
    selectedSize.value = sortedSizes.value[0];
  }
});

const { data: additiveCategories } = useQuery({
  queryKey: computed(() => ['kiosk-additive-categories', selectedSize.value]),
  queryFn: () => {
    if (selectedSize.value) {
      return storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id)
    }
  },
  initialData: [],
  enabled: computed(() => !!selectedSize.value),
});

const totalPrice = computed(() => {
  if (!selectedSize.value) return 0;
  const storePrice = selectedSize.value.storePrice;
  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => sum + add.storePrice, 0);
  return storePrice + additivePrice;
});

const onSizeSelect = (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value?.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = {};
};

const onAdditiveToggle = (category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[category.id] || [];
  const isSelected = current.some((a) => a.additiveId === additive.additiveId);

  if (category.isMultipleSelect) {
    selectedAdditives.value[category.id] = isSelected
      ? current.filter(a => a.additiveId !== additive.additiveId)
      : [...current, additive];
  } else {
    selectedAdditives.value[category.id] = isSelected ? [] : [additive];
  }
};

const isAdditiveSelected = (category: StoreAdditiveCategoryDTO, additiveId: number): boolean =>
  selectedAdditives.value[category.id]?.some((a) => a.additiveId === additiveId) || false;

const handleAddToCart = () => {
  if (!productDetails.value || !selectedSize.value) return;
  const allAdditives = Object.values(selectedAdditives.value).flat();
  cartStore.addToCart(productDetails.value, selectedSize.value, allAdditives, 1);
};
</script>

<style scoped></style>
