<template>
	<div class="relative bg-[#F3F4F9] pb-32 min-h-screen text-black">
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
			<div class="bg-white shadow-gray-200 shadow-xl px-8 pb-6 rounded-b-[48px] w-full">
				<div class="flex flex-col justify-center items-center">
					<LazyImage
						src="https://www.nicepng.com/png/full/106-1060376_starbucks-iced-coffee-png-vector-library-pumpkin-spice.png"
						alt="Изображение товара"
						class="w-32 h-52 object-contain"
					/>
					<p class="mt-7 font-semibold text-3xl">{{ productDetails.name }}</p>
					<p class="mt-1 text-slate-600 text-lg">{{ productDetails.description }}</p>
				</div>
				<div class="flex justify-between items-center gap-4 mt-12">
					<!-- Size Selection -->
					<div class="flex items-center gap-2 overflow-x-auto no-scrollbar">
						<KioskDetailsSizes
							v-for="size in productDetails.sizes"
							:key="size.id"
							:size="size"
							:is-selected="selectedSize?.id === size.id"
							@click:size="onSizeSelect"
						/>
					</div>
					<!-- Add to Cart Button -->
					<div class="flex items-center gap-6">
						<p class="text-3xl">
							{{ formatPrice(totalPrice) }}
						</p>
						<button
							@click="handleAddToCart"
							class="flex items-center gap-3 bg-primary p-5 rounded-full text-primary-foreground"
						>
							<Plus class="w-6 sm:w-8 h-6 sm:h-8" />
						</button>
					</div>
				</div>
			</div>
			<div class="mt-8">
				<KioskDetailsAdditivesSection
					:categories="additives"
					:isAdditiveSelected="isAdditiveSelected"
					@toggle-additive="onAdditiveToggle"
				/>
			</div>
			<!-- Additives Selection -->
		</div>
	</div>
</template>

<script setup lang="ts">
  import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductDetailsDTO, StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsLoading from '@/modules/kiosk/products/components/details/kiosk-details-loading.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import { Plus } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
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
  const additives = ref<StoreAdditiveCategoryDTO[]>([]);
  const selectedSize = ref<StoreProductSizeDetailsDTO | null>(null);
  const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>({});
  const quantity = ref<number>(1);
  const isLoading = ref<boolean>(true);
  const error = ref<string | null>(null);

  // Fetch product details based on productId prop
  const fetchProductDetails = async () => {
  isLoading.value = true;
  error.value = null;
  try {
  const details = await storeProductsService.getStoreProduct(props.productId);
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
  const fetchedAdditives = await storeAdditivesService.getStoreAdditiveCategories(sizeId);
  additives.value = fetchedAdditives;
  } catch (err) {
  console.error('Error fetching additives:', err);
  error.value = 'Failed to fetch additives.';
  }
  };
  // Compute total price
  const totalPrice = computed(() => {
  if (!selectedSize.value) return 0;
  const storePrice = selectedSize.value.storePrice;
  const additivePrice = Object.values(selectedAdditives.value)
  .flat()
  .reduce((sum, add) => sum + add.storePrice, 0);
  return (storePrice + additivePrice) * quantity.value;
  });

  // Handle size selection
  const onSizeSelect = async (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value?.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = {};
  await fetchAdditives(size.id);
  };
  // Toggle additive selection
  const onAdditiveToggle = (categoryId: number, additive: StoreAdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[categoryId] || [];
  const isSelected = current.some((a) => a.additiveId === additive.additiveId);
  if (isSelected) {
  selectedAdditives.value[categoryId] = current.filter((a) => a.additiveId !== additive.additiveId);
  } else {
  selectedAdditives.value[categoryId] = [...current, additive];
  }
  };
  // Check if additive is selected
  const isAdditiveSelected = (categoryId: number, additiveId: number): boolean =>
  selectedAdditives.value[categoryId]?.some((a) => a.additiveId === additiveId) || false;

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
