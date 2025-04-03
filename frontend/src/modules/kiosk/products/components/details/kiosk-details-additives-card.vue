<template>
	<div
		v-if="isVisible"
		@click="handleClick"
		:class="cardClasses"
		data-testid="additive-card"
	>
		<LazyImage
			:src="additive.imageUrl"
			alt="Изображение модификатора"
			:class="cn('rounded-[22px] size-16 md:size-24 object-contain', isPriceZero && 'md:size-16 rounded-[14px]')"
		/>

		<div
			:class="cn('flex flex-col justify-between items-start w-full h-full', isPriceZero && 'md:justify-center')"
		>
			<p
				class="text-lg md:text-2xl text-left line-clamp-2"
				data-testid="additive-name"
			>
				{{ additive.name }}
			</p>

			<div class="flex justify-between items-center gap-2 w-full">
				<p
					class="text-xl md:text-2xl text-left"
					:class="priceTextClass"
					data-testid="additive-price"
				>
					{{ priceDisplay }}
				</p>

				<template v-if="additive.isOutOfStock">
					<p class="bg-slate-200 px-3 py-1 rounded-xl font-medium text-slate-800 text-sm">
						Нет в наличии
					</p>
				</template>

				<template v-else>
					<button
						v-if="showSelectedIndicator"
						class="relative flex-shrink-0 rounded-full focus:outline-none size-8"
						:class="[isSelected || additive.isDefault ? 'bg-primary' : 'bg-gray-200']"
						data-testid="additive-button"
					>
						<span
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
import { cn } from '@/core/utils/tailwind.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { computed } from 'vue'

const props = defineProps<{
  additive: StoreAdditiveCategoryItemDTO;
  isSelected: boolean;
}>()

const emit = defineEmits(['click:additive'])

// Whether to render this card at all.
const isVisible = computed(() => {
  return !(props.additive.isDefault && props.additive.isHidden)
})

// Compute the base classes for the card.
const cardClasses = computed(() => {
  const classes = [
    'relative',
    'text-center',
    'bg-white',
    'rounded-[32px]',
    'p-6',
    'flex',
    'items-start',
    'border-2',
    'gap-4',
    'md:gap-6',
    'transition-all',
    'duration-300'
  ]

  if (props.additive.isOutOfStock) {
    classes.push('!cursor-not-allowed', 'opacity-60')
  } else {
    if (props.isSelected) {
      classes.push('bg-primary', 'border-primary', 'cursor-pointer')
    } else {
      classes.push('border-transparent', 'cursor-pointer')
    }
  }

  // Default additives have a fixed state.
  if (props.additive.isDefault) {
    classes.push('cursor-not-allowed', 'opacity-50', '!border-primary')
  }

  return classes
})

// Format the price display text.
const priceDisplay = computed(() => {
  if (props.additive.isDefault) {
    return 'Входит в состав'
  }
  if (props.additive.storePrice !== 0) {
    return formatPrice(props.additive.storePrice)
  }
  return '' // Hide text if price is 0.
})

// Text color class for the price.
const priceTextClass = computed(() => {
  return props.additive.isOutOfStock ? 'text-gray-400' : 'text-primary'
})

// Show the selected indicator only if selected (or default) and price is not 0.

const isPriceZero = computed(() => props.additive.storePrice === 0)

const showSelectedIndicator = computed(() => {
  return (props.isSelected || props.additive.isDefault) && !isPriceZero.value
})

const handleClick = () => {
  if (props.additive.isDefault || props.additive.isOutOfStock) {
    return
  }
  emit('click:additive', props.additive)
}
</script>

<style scoped></style>
