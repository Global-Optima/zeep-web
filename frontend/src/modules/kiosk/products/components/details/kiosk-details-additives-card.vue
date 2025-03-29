<template>
  <div
    v-if="!(additive.isDefault && additive.isHidden)"
    @click="handleClick"
    :class="[
    'relative text-center bg-white rounded-[24px] p-4 flex items-start border-2 gap-4 md:gap-6 transition-all duration-300',
    additive.isOutOfStock
      ? 'cursor-not-allowed opacity-60'
      : (isSelected ? 'bg-primary border-primary cursor-pointer' : 'border-transparent cursor-pointer'),
      additive.isDefault ? 'cursor-not-allowed opacity-50 !border-primary' : 'cursor-pointer',
  ]"
    data-testid="additive-card"
  >
    <LazyImage
      :src="additive.imageUrl"
      alt="Изображение добавки"
      class="rounded-md size-16 md:size-20 object-contain"
    />

    <div class="flex flex-col justify-between items-start w-full h-full">
      <p class="text-lg md:text-xl text-left line-clamp-2" data-testid="additive-name">
        {{ additive.name }}
      </p>

      <div class="flex justify-between items-center w-full gap-2 flex-wrap">
        <p
          class="text-left text-xl md:text-2xl"
          :class="additive.isOutOfStock ? 'text-gray-400' : 'text-primary'"
          data-testid="additive-price"
        >
          <span
            v-if="additive.isDefault"
            class="text-xl"
          >
						Входит в состав
					</span>

          <span v-else>{{formatPrice(additive.storePrice)}}</span>
        </p>


        <template v-if="additive.isOutOfStock">
          <p
            class="bg-slate-200 text-slate-800 text-sm font-medium px-3 py-1 rounded-xl"
          >
            Нет в наличии
          </p>
        </template>

        <template v-else>
          <button
            class="relative flex-shrink-0 rounded-full focus:outline-none size-8"
            :class="[isSelected || additive.isDefault ? 'bg-primary' : 'bg-gray-200']"
            data-testid="additive-button"
          >
            <span
              v-if="isSelected || additive.isDefault"
              class="absolute inset-0 bg-white m-auto rounded-full size-3"
              data-testid="additive-selected-indicator"
            ></span>
          </button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'

const {additive, isSelected} = defineProps<{
  additive: StoreAdditiveCategoryItemDTO;
  isSelected: boolean;
}>()

const emit = defineEmits(['click:additive'])

const handleClick = () => {
  if (additive.isDefault || additive.isOutOfStock) {
    return
  }
  emit('click:additive', additive)
}
</script>

<style scoped></style>
