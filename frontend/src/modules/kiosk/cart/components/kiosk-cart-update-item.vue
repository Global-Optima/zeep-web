<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent } from '@/core/components/ui/dialog'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsLoading from '@/modules/kiosk/products/components/details/kiosk-details-loading.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { ChevronLeft, Pencil } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'

// Props & Emits
const props = defineProps<{
  isOpen: boolean;
  cartItem: CartItem;
}>();

const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'update', updatedSize: StoreProductSizeDetailsDTO, updatedAdditives: StoreAdditiveCategoryItemDTO[]): void;
}>();

// Individual state variables
const productDetails = ref(props.cartItem.product);
const selectedSize = ref<StoreProductSizeDetailsDTO>(props.cartItem.size);
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>({});
const additivesCategories = ref<StoreAdditiveCategoryDTO[]>([]);
const isLoading = ref(false);
const errorMessage = ref<string | null>(null);

// Fetch additives for the selected size
const fetchAdditives = async () => {
  try {
    isLoading.value = true;
    const fetchedCategories = await storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id);

    // Merge fetched additives with pre-selected ones
    additivesCategories.value = fetchedCategories.map(category => ({
      ...category,
      additives: category.additives.map(additive => ({
        ...additive,
        isSelected: isAdditiveSelected(category, additive.additiveId),
      })),
    }));
  } catch {
    errorMessage.value = 'Ошибка при загрузке добавок';
  } finally {
    isLoading.value = false;
  }
};

// Initialize state & fetch additives
const resetAndFetchAdditives = async () => {
  selectedSize.value = props.cartItem.size;
  selectedAdditives.value = props.cartItem.additives.reduce<Record<number, StoreAdditiveCategoryItemDTO[]>>((acc, additive) => {
    acc[additive.categoryId] = acc[additive.categoryId] || [];
    acc[additive.categoryId].push(additive);
    return acc;
  }, {});
  await fetchAdditives();
};

// Preserve additives across size changes
const preserveSelectedAdditives = () => {
  const previousAdditives = { ...selectedAdditives.value };
  Object.keys(previousAdditives).forEach(categoryId => {
    const category = additivesCategories.value.find(cat => cat.id === Number(categoryId));
    if (category) {
      selectedAdditives.value[category.id] = previousAdditives[category.id] || [];
    }
  });
};

// Handle size selection
const onSizeSelect = async (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = {};
  await fetchAdditives();
  preserveSelectedAdditives();
};

// Computed total price
const totalPrice = computed(() => {
  const storePrice = selectedSize.value.storePrice;
  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => sum + add.storePrice, 0);
  return storePrice + additivePrice;
});

// Check if additive is selected
const isAdditiveSelected = (category: StoreAdditiveCategoryDTO, additiveId: number) =>
  selectedAdditives.value[category.id]?.some(a => a.additiveId === additiveId) || false;

// Handle additive selection
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


const sortedSizes = computed(() =>
  productDetails.value?.sizes ? [...productDetails.value.sizes].sort((a, b) => a.size - b.size) : []
);

// Watch for dialog open state and reinitialize
watch(
  () => props.isOpen,
  async open => {
    if (open) await resetAndFetchAdditives();
  }
);

// On mounted, check and fetch data if dialog is already open
onMounted(async () => {
  if (props.isOpen) await resetAndFetchAdditives();
});

// Handle update action
const handleUpdate = () => {
  const updatedAdditives = Object.values(selectedAdditives.value).flat();
  emits('update', selectedSize.value, updatedAdditives);
  emits('close');
};
</script>

<template>
	<Dialog
		:open="isOpen"
		@update:open="emits('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="p-0 border-none !rounded-[52px] max-w-[80vw] h-[90vh] overflow-clip"
		>
			<div class="relative bg-[#F3F4F9] overflow-y-auto text-black no-scrollbar">
				<!-- Loading State -->
				<KioskDetailsLoading v-if="isLoading" />
				<!-- Error State -->
				<div
					v-else-if="errorMessage"
					class="p-4 text-red-500"
				>
					{{ errorMessage }}
				</div>
				<!-- Product Content -->
				<div
					v-else
					class="pb-44"
				>
					<div class="bg-white shadow-gray-200 shadow-xl px-8 pb-6 rounded-b-[48px] w-full">
						<header class="flex justify-between items-center gap-6 pt-6">
							<Button
								size="icon"
								variant="ghost"
								class="size-14"
								@click="emits('close')"
							>
								<ChevronLeft
									class="size-14 text-gray-400"
									stroke-width="1.6"
								/>
							</Button>
							<KioskProductRecipeDialog :nutrition="selectedSize.totalNutrition" />
						</header>

						<div class="flex flex-col justify-center items-center">
							<LazyImage
								:src="productDetails.imageUrl"
								alt="Изображение товара"
								class="w-38 h-64 object-contain"
							/>
							<p class="mt-7 font-semibold text-4xl">{{ productDetails.name }}</p>
							<p class="mt-2 text-slate-600 text-xl">{{ productDetails.description }}</p>
						</div>

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
										@click="handleUpdate"
										class="flex items-center gap-3 bg-primary p-5 rounded-full text-primary-foreground"
									>
										<Pencil class="size-8 sm:size-12" />
									</button>
								</div>
							</div>
						</div>
					</div>

					<!-- Additives Selection -->
					<div class="mt-10">
						<KioskDetailsAdditivesSection
							:categories="additivesCategories"
							:isAdditiveSelected="isAdditiveSelected"
							@toggle-additive="onAdditiveToggle"
						/>
					</div>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
