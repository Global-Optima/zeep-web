<template>
	<div
		@click="handleClick"
		:class="[
      'text-center bg-white rounded-3xl p-5 min-w-44 max-w-44 flex flex-col justify-between',
      isSelected ? 'bg-primary border-primary border-2' : 'border-2 border-transparent',
      isDefaultAdditive ? 'cursor-not-allowed ' : 'cursor-pointer',
    ]"
		data-testid="additive-card"
	>
		<img
			:src="additive.imageUrl"
			alt="Additive Image"
			class="w-full h-24 object-contain"
			data-testid="additive-image"
		/>
		<p
			class="mt-3 flex-grow text-sm sm:text-base"
			data-testid="additive-name"
		>
			{{ additive.name }}
		</p>

		<div class="flex items-center justify-between mt-5">
			<p
				v-if="!isDefaultAdditive"
				:class="[
          'text-gray-400 text-lg sm:text-xl',
          isSelected ? 'text-black' : '',
        ]"
				data-testid="additive-price"
			>
				{{formatPrice(additive.price) }}
			</p>
			<button
				class="relative h-6 w-6 sm:w-8 sm:h-8 rounded-full focus:outline-none"
				:class="isSelected ? 'bg-primary' : 'bg-gray-200'"
				:disabled="isDefaultAdditive"
				data-testid="additive-button"
			>
				<span
					v-if="isSelected"
					class="absolute inset-0 m-auto sm:w-4 sm:h-4 h-3 w-3 rounded-full bg-white"
					data-testid="additive-selected-indicator"
				></span>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { Additive } from '@/modules/kiosk/products/models/product.model'
import { computed } from 'vue'

const props = defineProps<{
  additive: Additive;
  selectedAdditives: Additive[];
  defaultAdditives: Additive[];
}>();

const emit = defineEmits(['click:additive']);

const { additive, selectedAdditives, defaultAdditives } = props;

const isDefaultAdditive = computed(() => {
  return defaultAdditives.some((item) => item.id === additive.id);
});

const isSelected = computed(() => {
  return selectedAdditives.some((item) => item.id === additive.id) || isDefaultAdditive.value;
});

const handleClick = () => {
  if (isDefaultAdditive.value) {
    return;
  }
  emit('click:additive', additive);
};
</script>

<style scoped></style>
