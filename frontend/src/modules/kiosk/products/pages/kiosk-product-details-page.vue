<template>
	<div class="relative bg-gray-100 min-h-screen text-black">
		<!-- Back Button -->
		<KioskDetailsBack />

		<!-- Loading State -->
		<KioskDetailsLoading v-if="state.isLoading" />

		<!-- Product Content -->
		<div
			v-else-if="state.productDetails"
			class="pb-44"
		>
			<!-- Product Image -->
			<KioskDetailsProductImage
				:imageUrl="state.productDetails.imageUrl"
				:altText="state.productDetails.name"
			/>

			<!-- Product Information -->
			<KioskDetailsProductInfo
				:name="state.productDetails.name"
				:description="state.productDetails.description"
				:energy="calculatedEnergy"
			/>

			<!-- Additives Selection -->
			<KioskDetailsAdditivesSection
				:categories="sortedAdditiveCategories"
				:isAdditiveDefault="isAdditiveDefault"
				:isAdditiveSelected="isAdditiveSelected"
				@toggle-additive="onAdditiveToggle"
			/>
		</div>

		<!-- Fixed Bottom Section -->
		<KioskDetailsBottomBar
			v-if="state.productDetails"
			:sizes="state.productDetails.sizes"
			:selectedSizeId="state.selectedSize?.id"
			:totalPrice="totalPrice"
			:formatPrice="formatPrice"
			@select-size="onSizeSelect"
			@addToCart="handleAddToCart"
		/>
	</div>
</template>

<script setup lang="ts">
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { productService } from '@/modules/kiosk/products/services/products.service'
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { formatPrice } from '@/core/utils/price.utils'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsBack from '@/modules/kiosk/products/components/details/kiosk-details-back.vue'
import KioskDetailsBottomBar from '@/modules/kiosk/products/components/details/kiosk-details-bottom-bar.vue'
import KioskDetailsLoading from '@/modules/kiosk/products/components/details/kiosk-details-loading.vue'
import KioskDetailsProductImage from '@/modules/kiosk/products/components/details/kiosk-details-product-image.vue'
import KioskDetailsProductInfo from '@/modules/kiosk/products/components/details/kiosk-details-product-info.vue'
import type { AdditiveCategoryDTO, AdditiveDTO, ProductSizeDTO, StoreProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();

const state = reactive({
  productId: Number(route.params.id),
  productDetails: null as StoreProductDetailsDTO | null,
  additives: [] as AdditiveCategoryDTO[],
  selectedSize: null as ProductSizeDTO | null,
  selectedAdditives: {} as Record<number, AdditiveDTO[]>,
  quantity: 1,
  isLoading: true,
  error: null as string | null,
});

const fetchProductDetails = async () => {
  try {
    state.isLoading = true;
    const productDetails = await productService.getStoreProductDetails(state.productId);
    state.productDetails = productDetails;

    if (productDetails.sizes.length > 0) {
      state.selectedSize = productDetails.sizes[0];
      await fetchAdditives(state.selectedSize.id);
    }
  } catch {
    state.error = 'Failed to load product details.';
  } finally {
    state.isLoading = false;
  }
};

const fetchAdditives = async (sizeId: number) => {
  try {
    state.additives = await productService.getAdditiveCategoriesByProductSize(sizeId);
  } catch {
    state.error = 'Failed to fetch additives.';
  }
};

const totalPrice = computed(() => {
  const basePrice = state.selectedSize?.basePrice || 0;
  const additivePrice = Object.values(state.selectedAdditives)
    .flat()
    .reduce((sum, add) => sum + add.price, 0);
  return (basePrice + additivePrice) * state.quantity;
});

const calculatedEnergy = computed(() => {
  const details = state.productDetails;
  return details
    ? { ccal: 400, proteins: 20, carbs: 13, fats: 10 } // Placeholder for actual logic
    : { ccal: 0, proteins: 0, carbs: 0, fats: 0 };
});

const sortedAdditiveCategories = computed(() =>
  state.additives.map((category) => ({
    ...category,
    additives: [
      ...category.additives.filter((a) => isAdditiveDefault(a.id)),
      ...category.additives.filter((a) => !isAdditiveDefault(a.id)),
    ],
  }))
);

const onSizeSelect = (size: ProductSizeDTO) => {
  if (state.selectedSize?.id === size.id) return;
  state.selectedSize = size;
  state.selectedAdditives = {};
  fetchAdditives(size.id);
};

const onAdditiveToggle = (categoryId: number, additive: AdditiveDTO) => {
  const current = state.selectedAdditives[categoryId] || [];
  const isSelected = current.some((a) => a.id === additive.id);
  state.selectedAdditives[categoryId] = isSelected
    ? current.filter((a) => a.id !== additive.id)
    : [...current, additive];
};

const isAdditiveSelected = (categoryId: number, additiveId: number) =>
  state.selectedAdditives[categoryId]?.some((a) => a.id === additiveId) || false;

const isAdditiveDefault = (additiveId: number) =>
  state.productDetails?.defaultAdditives.some((add) => add.id === additiveId) || false;

const handleAddToCart = () => {
  if (!state.productDetails || !state.selectedSize) return;
  const allAdditives = Object.values(state.selectedAdditives).flat();
  cartStore.addToCart(state.productDetails, state.selectedSize, allAdditives, state.quantity);
  router.back();
};

onMounted(fetchProductDetails);
</script>

<style scoped lang="scss"></style>
