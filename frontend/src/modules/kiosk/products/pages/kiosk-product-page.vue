<template>
	<div class="relative bg-[#F3F4F9] pb-32 text-black no-scrollbar">
		<!-- Loading State -->
		<KioskDetailsLoading v-if="isFetching" />
		<!-- Error State -->
		<div
			v-else-if="isError"
			class="p-4 text-red-500"
		>
			{{ error?.message || 'Ошибка загрузки данных' }}
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

					<KioskProductRecipeDialog :product="productDetails" />
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
					:categories="additives ?? []"
					:isAdditiveSelected="isAdditiveSelected"
					@toggle-additive="onAdditiveToggle"
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type {
  StoreAdditiveCategoryItemDTO
} from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsLoading from '@/modules/kiosk/products/components/details/kiosk-details-loading.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, Plus } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();

// ✅ Product ID
const productId = computed(() => Number(route.params.id));

// ✅ Go back to Kiosk Home
const onBackClick = () => router.push({ name: getRouteName('KIOSK_HOME') });

// ✅ Selected Size & Additives
const selectedSize = ref<StoreProductSizeDetailsDTO | null>(null);
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>({});

// ✅ Fetch Product Details with `useQuery`
const { data: productDetails, isLoading: isFetching, isError, error } = useQuery({
  queryKey: ['kiosk-product-details', productId],
  queryFn: () => storeProductsService.getStoreProduct(productId.value),
  enabled: productId.value > 0,
});

const sortedSizes = computed(() =>
  productDetails.value?.sizes ? [...productDetails.value.sizes].sort((a, b) => a.size - b.size) : []
);

watch(sortedSizes, (sizes) => {
  if (sizes.length > 0) {
    selectedSize.value = sizes[0];
  }
});

// ✅ Fetch Additives with `useQuery`
const { data: additives } = useQuery({
  queryKey: ['kiosk-additive-categories', productId],
  queryFn: () => {
    if (selectedSize.value) {
      return storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id)
    }
  },
  initialData: [],
  enabled: computed(() => !!selectedSize.value), // Runs when size is selected
});

// ✅ Compute Total Price
const totalPrice = computed(() => {
  if (!selectedSize.value) return 0;
  const storePrice = selectedSize.value.storePrice;
  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => sum + add.storePrice, 0);
  return storePrice + additivePrice;
});

// ✅ Handle Size Selection
const onSizeSelect = (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value?.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = {};
};

// ✅ Toggle Additive Selection
const onAdditiveToggle = (categoryId: number, additive: StoreAdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[categoryId] || [];
  const isSelected = current.some((a) => a.additiveId === additive.additiveId);
  selectedAdditives.value[categoryId] = isSelected
    ? current.filter((a) => a.additiveId !== additive.additiveId)
    : [...current, additive];
};

// ✅ Check if Additive is Selected
const isAdditiveSelected = (categoryId: number, additiveId: number): boolean =>
  selectedAdditives.value[categoryId]?.some((a) => a.additiveId === additiveId) || false;

// ✅ Handle Add to Cart
const handleAddToCart = () => {
  if (!productDetails.value || !selectedSize.value) return;
  const allAdditives = Object.values(selectedAdditives.value).flat();
  cartStore.addToCart(productDetails.value, selectedSize.value, allAdditives, 1);
};
</script>

<style scoped></style>
