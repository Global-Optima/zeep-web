<template>
	<div class="relative bg-gray-100 pb-32 min-h-screen text-black">
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
			@select-size="onSizeSelect"
			@addToCart="handleAddToCart"
		/>
	</div>
</template>

<script setup lang="ts">
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { productService } from '@/modules/kiosk/products/services/products.service'
import { computed, onMounted, reactive, watch } from 'vue'

import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsBottomBar from '@/modules/kiosk/products/components/details/kiosk-details-bottom-bar.vue'
import KioskDetailsLoading from '@/modules/kiosk/products/components/details/kiosk-details-loading.vue'
import KioskDetailsProductImage from '@/modules/kiosk/products/components/details/kiosk-details-product-image.vue'
import KioskDetailsProductInfo from '@/modules/kiosk/products/components/details/kiosk-details-product-info.vue'

import type {
  AdditiveCategoryDTO,
  AdditiveDTO,
  ProductSizeDTO,
  StoreProductDetailsDTO,
} from '@/modules/kiosk/products/models/product.model'

// Define props
const props = defineProps<{
  productId: number;
}>();

// Define emitted events
const emit = defineEmits(['close']);

const cartStore = useCartStore();

const state = reactive({
  productDetails: null as StoreProductDetailsDTO | null,
  additives: [] as AdditiveCategoryDTO[],
  selectedSize: null as ProductSizeDTO | null,
  selectedAdditives: {} as Record<number, AdditiveDTO[]>,
  quantity: 1,
  isLoading: true,
  error: null as string | null,
});

// Fetch product details based on productId prop
const fetchProductDetails = async () => {
  try {
    state.isLoading = true;
    const productDetails = await productService.getStoreProductDetails(props.productId);
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

// Fetch additives based on selected size
const fetchAdditives = async (sizeId: number) => {
  try {
    state.additives = await productService.getAdditiveCategoriesByProductSize(sizeId);
  } catch {
    state.error = 'Failed to fetch additives.';
  }
};

// Compute total price
const totalPrice = computed(() => {
  const basePrice = state.selectedSize?.basePrice || 0;
  const additivePrice = Object.values(state.selectedAdditives)
    .flat()
    .reduce((sum, add) => sum + add.price, 0);
  return (basePrice + additivePrice) * state.quantity;
});

// Placeholder for energy calculation logic
const calculatedEnergy = computed(() => {
  const details = state.productDetails;
  return details
    ? { ccal: 400, proteins: 20, carbs: 13, fats: 10 }
    : { ccal: 0, proteins: 0, carbs: 0, fats: 0 };
});

// Sort additives categories
const sortedAdditiveCategories = computed(() =>
  state.additives.map((category) => ({
    ...category,
    additives: [
      ...category.additives.filter((a) => isAdditiveDefault(a.id)),
      ...category.additives.filter((a) => !isAdditiveDefault(a.id)),
    ],
  }))
);

// Handle size selection
const onSizeSelect = (size: ProductSizeDTO) => {
  if (state.selectedSize?.id === size.id) return;
  state.selectedSize = size;
  state.selectedAdditives = {};
  fetchAdditives(size.id);
};

// Toggle additive selection
const onAdditiveToggle = (categoryId: number, additive: AdditiveDTO) => {
  const current = state.selectedAdditives[categoryId] || [];
  const isSelected = current.some((a) => a.id === additive.id);
  state.selectedAdditives[categoryId] = isSelected
    ? current.filter((a) => a.id !== additive.id)
    : [...current, additive];
};

// Check if additive is selected
const isAdditiveSelected = (categoryId: number, additiveId: number) =>
  state.selectedAdditives[categoryId]?.some((a) => a.id === additiveId) || false;

// Check if additive is default
const isAdditiveDefault = (additiveId: number) =>
  state.productDetails?.defaultAdditives.some((add) => add.id === additiveId) || false;

// Handle add to cart action
const handleAddToCart = () => {
  if (!state.productDetails || !state.selectedSize) return;
  const allAdditives = Object.values(state.selectedAdditives).flat();
  cartStore.addToCart(state.productDetails, state.selectedSize, allAdditives, state.quantity);
  // Emit close event to parent component
  emit('close');
};

// Fetch product details when component is mounted
onMounted(fetchProductDetails);

// Watch for changes in productId prop
watch(
  () => props.productId,
  (newProductId, oldProductId) => {
    if (newProductId !== oldProductId) {

      state.productDetails = null;
      state.additives = [];
      state.selectedSize = null;
      state.selectedAdditives = {};
      state.quantity = 1;
      state.isLoading = true;
      state.error = null;

      fetchProductDetails();
    }
  }
);
</script>

<style scoped lang="scss">
/* Add your styles here */
</style>
