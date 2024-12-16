<template>
	<div class="relative bg-gray-100 pb-32 min-h-screen text-black">
		<!-- Loading State -->
		<KioskDetailsLoading v-if="isLoading" />

		<!-- Error State -->
		<div
			v-else-if="error"
			class="p-4 text-red-500"
		>
			{{ error }}
		</div>

		<!-- Product Content -->
		<div
			v-else-if="productDetails"
			class="pb-44"
		>
			<!-- Product Image -->
			<KioskDetailsProductImage
				:imageUrl="productDetails.imageUrl"
				:altText="productDetails.name"
			/>

			<!-- Product Information -->
			<KioskDetailsProductInfo
				:name="productDetails.name"
				:description="productDetails.description"
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
			v-if="productDetails"
			:sizes="productDetails.sizes"
			:selectedSizeId="selectedSize?.id"
			:totalPrice="totalPrice"
			@select-size="onSizeSelect"
			@addToCart="handleAddToCart"
		/>
	</div>
</template>

<script setup lang="ts">
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { computed, defineProps, onMounted, ref, watch } from 'vue'

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
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'

// Define props
const props = defineProps<{
  productId: number;
}>();

// Define emitted events
const emit = defineEmits<{
  (e: 'close'): void;
}>();

const cartStore = useCartStore();

// Reactive state
const productDetails = ref<StoreProductDetailsDTO | null>(null);
const additives = ref<AdditiveCategoryDTO[]>([]);
const selectedSize = ref<ProductSizeDTO | null>(null);
const selectedAdditives = ref<Record<number, AdditiveDTO[]>>({});
const quantity = ref<number>(1);
const isLoading = ref<boolean>(true);
const error = ref<string | null>(null);

const {currentStoreId} = useCurrentStoreStore()

// Fetch product details based on productId prop
const fetchProductDetails = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    if (!currentStoreId) return

    const details = await productsService.getStoreProductDetails(props.productId, currentStoreId);
    productDetails.value = details;

    if (details.sizes.length > 0) {
      selectedSize.value = details.sizes[0];
      await fetchAdditives(details.sizes[0].id);
    }
  } catch (err) {
    console.error('Error fetching product details:', err);
    error.value = 'Failed to load product details.';
  } finally {
    isLoading.value = false;
  }
};

// Fetch additives based on selected size
const fetchAdditives = async (sizeId: number) => {
  try {
    const fetchedAdditives = await productsService.getAdditiveCategoriesByProductSize(sizeId);
    additives.value = fetchedAdditives;
  } catch (err) {
    console.error('Error fetching additives:', err);
    error.value = 'Failed to fetch additives.';
  }
};

// Compute total price
const totalPrice = computed(() => {
  if (!selectedSize.value) return 0;
  const basePrice = selectedSize.value.basePrice;
  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => sum + add.price, 0);
  return (basePrice + additivePrice) * quantity.value;
});

// Placeholder for energy calculation logic
const calculatedEnergy = computed(() => {
  if (!productDetails.value) {
    return { ccal: 0, proteins: 0, carbs: 0, fats: 0 };
  }
  // Replace with actual calculation logic
  return { ccal: 400, proteins: 20, carbs: 13, fats: 10 };
});

// Sort additives categories with memoization
const sortedAdditiveCategories = computed(() =>
  additives.value.map((category) => ({
    ...category,
    additives: [
      ...category.additives.filter((a) => isAdditiveDefault(a.id)),
      ...category.additives.filter((a) => !isAdditiveDefault(a.id)),
    ],
  }))
);

// Handle size selection
const onSizeSelect = async (size: ProductSizeDTO) => {
  if (selectedSize.value?.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = {};
  await fetchAdditives(size.id);
};

// Toggle additive selection
const onAdditiveToggle = (categoryId: number, additive: AdditiveDTO) => {
  const current = selectedAdditives.value[categoryId] || [];
  const isSelected = current.some((a) => a.id === additive.id);
  if (isSelected) {
    selectedAdditives.value[categoryId] = current.filter((a) => a.id !== additive.id);
  } else {
    selectedAdditives.value[categoryId] = [...current, additive];
  }
};

// Check if additive is selected
const isAdditiveSelected = (categoryId: number, additiveId: number): boolean =>
  selectedAdditives.value[categoryId]?.some((a) => a.id === additiveId) || false;

// Check if additive is default
const isAdditiveDefault = (additiveId: number): boolean =>
  productDetails.value?.defaultAdditives.some((add) => add.id === additiveId) || false;

// Handle add to cart action
const handleAddToCart = () => {
  if (!productDetails.value || !selectedSize.value) return;
  const allAdditives = Object.values(selectedAdditives.value).flat();
  cartStore.addToCart(productDetails.value, selectedSize.value, allAdditives, quantity.value);
  // Emit close event to parent component
  emit('close');
};

// Fetch product details when component is mounted
onMounted(fetchProductDetails);

// Watch for changes in productId prop
watch(
  () => props.productId,
  async (newProductId, oldProductId) => {
    if (newProductId !== oldProductId) {
      // Reset state
      productDetails.value = null;
      additives.value = [];
      selectedSize.value = null;
      selectedAdditives.value = {};
      quantity.value = 1;
      error.value = null;
      await fetchProductDetails();
    }
  }
);
</script>

<style scoped lang="scss">
/* Add your styles here */
</style>
