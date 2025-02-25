<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { computed, onMounted, ref } from "vue"

const props = defineProps<{
  src?: string;
  alt?: string;
  class?: string;
}>();

// Reactive state
const imageSrc = ref(props.src);
const isLoading = ref(true);
const hasError = ref(false);

// Resolve the fallback image correctly
const fallbackImage = new URL("/src/core/assets/images/placeholder.svg", import.meta.url).href;

onMounted(() => {
  if (!props.src) {
    useFallback();
    return;
  }
  const img = new Image();
  img.src = props.src;
  img.onload = () => {
    isLoading.value = false;
  };
  img.onerror = () => {
    useFallback();
  };
});

// Fallback logic
const useFallback = () => {
  hasError.value = true;
  isLoading.value = false;
  imageSrc.value = fallbackImage;
};

// Computed classes for main image (smooth fade-in)
const imageClass = computed(() =>
  cn(
    "w-full h-full object-cover transition-opacity duration-500",
    props.class,
    { "opacity-0": isLoading.value }
  )
);
</script>

<template>
	<!-- Skeleton Loader -->
	<div
		v-if="isLoading"
		:class="cn('w-full h-full', props.class, 'bg-gray-50 animate-pulse duration-500')"
	></div>

	<!-- Main Image -->
	<img
		v-else-if="!hasError"
		:src="imageSrc"
		:alt="alt || 'Image'"
		:class="imageClass"
		loading="lazy"
	/>

	<!-- Fallback Image (only shown when error occurs) -->
	<img
		v-else
		:src="fallbackImage"
		:alt="alt || 'Placeholder'"
		:class="cn('w-full h-full', props.class, 'object-cover')"
	/>
</template>
