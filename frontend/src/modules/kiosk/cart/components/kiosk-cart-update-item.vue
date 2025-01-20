<script setup lang="ts">
import {
  Dialog,
  DialogContent
} from '@/core/components/ui/dialog'
import type { AdditiveCategoryItemDTO } from '@/modules/admin/additives/models/additives.model'
import type { StoreAdditiveCategoryDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsBottomBar from '@/modules/kiosk/products/components/details/kiosk-details-bottom-bar.vue'
import KioskDetailsProductImage from '@/modules/kiosk/products/components/details/kiosk-details-product-image.vue'
import KioskDetailsProductInfo from '@/modules/kiosk/products/components/details/kiosk-details-product-info.vue'
import { X } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'

interface UpdateCartDialogProps {
  isOpen: boolean;
  cartItem: CartItem;
  onClose: () => void;
  onUpdate: (updatedSize: StoreProductSizeDTO, updatedAdditives: AdditiveCategoryItemDTO[]) => void;
}

const { isOpen, cartItem, onClose, onUpdate } = defineProps<UpdateCartDialogProps>();

// Individual state variables
const productDetails = ref(cartItem.product);
const selectedSize = ref(cartItem.size);
const selectedAdditives = ref(
  cartItem.additives.reduce<Record<number, AdditiveCategoryItemDTO[]>>((acc, additive) => {
    acc[additive.categoryId] = acc[additive.categoryId] || [];
    acc[additive.categoryId].push(additive);
    return acc;
  }, {})
);
const additivesCategories = ref<StoreAdditiveCategoryDTO[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

// Fetch additives for the selected size
const fetchAdditives = async () => {
  try {
    isLoading.value = true;
    const fetchedCategories = await storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id);

    // Merge fetched additives with pre-selected ones
    additivesCategories.value = fetchedCategories.map((category) => ({
      ...category,
      additives: category.additives.map((additive) => ({
        ...additive,
        isSelected: isAdditiveSelected(category.id, additive.id),
      })),
    }));
  } catch {
    error.value = 'Failed to load additives.';
  } finally {
    isLoading.value = false;
  }
};

// Reset dialog state and fetch additives
const resetAndFetchAdditives = async () => {
  selectedSize.value = cartItem.size;
  selectedAdditives.value = cartItem.additives.reduce<Record<number, AdditiveCategoryItemDTO[]>>((acc, additive) => {
    acc[additive.categoryId] = acc[additive.categoryId] || [];
    acc[additive.categoryId].push(additive);
    return acc;
  }, {});
  await fetchAdditives();
};

// Preserve additives across size changes
const preserveSelectedAdditives = () => {
  const previousAdditives = { ...selectedAdditives.value };
  Object.keys(previousAdditives).forEach((categoryId) => {
    const category = additivesCategories.value.find((cat) => cat.id === Number(categoryId));
    if (category) {
      selectedAdditives.value[category.id] = previousAdditives[category.id] || [];
    }
  });
};

// Handle size selection
const onSizeSelect = async (size: StoreProductSizeDTO) => {
  if (selectedSize.value?.id === size.id) return;
  selectedSize.value = size;
  selectedAdditives.value = []
  await fetchAdditives();
  preserveSelectedAdditives();
};


// Computed total price
const totalPrice = computed(() => {
  const basePrice = selectedSize.value?.basePrice || 0;
  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => sum + add.price, 0);
  return basePrice + additivePrice;
});

// Computed energy
const calculatedEnergy = computed(() => {
  const baseEnergy = { ccal: 100, proteins: 5, carbs: 20, fats: 2 };
  const additiveEnergy = Object.values(selectedAdditives.value)
    .flat()
    .reduce(
      (total) => ({
        ccal: total.ccal + 200,
        proteins: total.proteins + 20,
        carbs: total.carbs + 10,
        fats: total.fats + 5,
      }),
      { ccal: 0, proteins: 0, carbs: 0, fats: 0 }
    );

  return {
    ccal: baseEnergy.ccal + additiveEnergy.ccal,
    proteins: baseEnergy.proteins + additiveEnergy.proteins,
    carbs: baseEnergy.carbs + additiveEnergy.carbs,
    fats: baseEnergy.fats + additiveEnergy.fats,
  };
});

// Check if additive is selected
const isAdditiveSelected = (categoryId: number, additiveId: number) =>
  selectedAdditives.value[categoryId]?.some((a) => a.id === additiveId) || false;

// Handle additive selection
const onAdditiveToggle = (categoryId: number, additive: AdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[categoryId] || [];
  const isSelected = current.some((a) => a.id === additive.id);
  selectedAdditives.value[categoryId] = isSelected
    ? current.filter((a) => a.id !== additive.id)
    : [...current, additive];
};

// Watch for dialog open state and reinitialize
watch(
  () => isOpen,
  async (open) => {
    if (open) {
      await resetAndFetchAdditives();
    }
  }
);

// On mounted, check and fetch data if dialog is already open
onMounted(async () => {
  if (isOpen) {
    await resetAndFetchAdditives();
  }
});

// Handle update action
const handleUpdate = () => {
  const updatedAdditives = Object.values(selectedAdditives.value).flat();
  onUpdate(selectedSize.value, updatedAdditives);
  onClose();
};
</script>

<template>
	<Dialog
		:open="isOpen"
		@update:open="onClose"
	>
		<DialogContent
			:include-close-button="false"
			class="bg-[#F5F5F7] p-0 border-none !rounded-[52px] max-w-[80vw] h-[90vh] overflow-clip"
		>
			<div class="relative bg-gray-100 pb-44 w-full text-black overflow-y-auto no-scrollbar">
				<button
					@click="onClose"
					class="top-14 right-14 z-10 fixed bg-gray-300 bg-opacity-90 p-3 rounded-full"
				>
					<X class="w-8 h-8 text-gray-800" />
				</button>

				<!-- Product Image -->
				<KioskDetailsProductImage
					:imageUrl="productDetails.imageUrl"
					:altText="productDetails.name"
					image-class="h-[500px]"
					overlay-class="h-[500px]"
				/>

				<!-- Product Information -->
				<KioskDetailsProductInfo
					v-if="selectedSize"
					:name="productDetails.name"
					:description="productDetails.description"
					:energy="calculatedEnergy"
					:ingredients="selectedSize?.ingredients"
				/>

				<!-- Additives Selection -->
				<KioskDetailsAdditivesSection
					:categories="additivesCategories"
					:isAdditiveSelected="isAdditiveSelected"
					@toggle-additive="onAdditiveToggle"
				/>

				<!-- Sizes and Price -->
				<KioskDetailsBottomBar
					:sizes="productDetails.sizes"
					:selectedSizeId="selectedSize.id"
					:totalPrice="totalPrice"
					display-icon="update"
					@select-size="onSizeSelect"
					@add-to-cart="handleUpdate"
				/>
			</div>
		</DialogContent>
	</Dialog>
</template>

<style scoped>
/* Add any custom styles here */
</style>
