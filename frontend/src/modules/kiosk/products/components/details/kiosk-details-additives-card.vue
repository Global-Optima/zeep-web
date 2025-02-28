<template>
	<div
		@click="handleClick"
		:class="[
      'text-center bg-white rounded-[32px] p-5 flex items-start border-2 gap-8',
      isSelected ? 'bg-primary border-primary' : 'border-transparent',
      isDefault  ? 'cursor-not-allowed opacity-60 !border-primary' : 'cursor-pointer',

    ]"
		data-testid="additive-card"
	>
		<LazyImage
			:src="additive.imageUrl"
			alt="Изображение добавки"
			class="rounded-md w-24 h-24 object-contain"
		/>

		<div class="flex flex-col justify-between items-start w-full h-full">
			<p
				class="text-xl text-left line-clamp-2"
				data-testid="additive-name"
			>
				{{ additive.name }}
			</p>

			<div class="flex justify-between items-center w-full">
				<p
					:class="[
          'text-primary text-2xl',
          isSelected ? 'text-black' : '',
        ]"
					data-testid="additive-price"
				>
					{{ isDefault ? formatPrice(0) : formatPrice(additive.storePrice) }}
				</p>
				<button
					class="relative rounded-full focus:outline-none sm:size-8"
					:class="[
          isSelected ? 'bg-primary' : 'bg-gray-200',
          isDefault ? 'cursor-not-allowed opacity-50' : '',
        ]"
					:disabled="isDefault"
					data-testid="additive-button"
				>
					<span
						v-if="isSelected"
						class="absolute inset-0 bg-white m-auto rounded-full w-3 sm:w-4 h-3 sm:h-4"
						data-testid="additive-selected-indicator"
					></span>
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'

const {additive, isSelected, isDefault} = defineProps<{
  additive: StoreAdditiveCategoryItemDTO;
  isSelected: boolean;
  isDefault: boolean;
}>()

const emit = defineEmits(['click:additive'])

const handleClick = () => {
  if (isDefault) {
    return
  }
  emit('click:additive', additive)
}
</script>

<style scoped></style>
