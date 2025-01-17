<template>
	<section
		class="top-0 z-10 sticky flex items-center gap-2 px-4 py-4 sm:py-6 w-full overflow-x-auto no-scrollbar"
		data-testid="search-section"
	>
		<div
			class="relative flex-grow"
			data-testid="search-input-container"
		>
			<input
				type="text"
				:value="searchTerm"
				@input="updateSearchTerm"
				@focus="handleFocus"
				@blur="handleBlur"
				:class="[
          'text-base sm:text-xl rounded-full p-5 sm:px-8 sm:py-6 focus:outline-none transition-all font-medium',
          {
            'w-28 sm:w-36 text-center': !isInputFocused,
            'w-56 sm:w-80': isInputFocused,
          }
        ]"
				placeholder="Поиск"
				data-testid="search-input"
			/>

			<button
				v-if="searchTerm"
				@click="clearSearchTerm"
				class="top-5 sm:top-6 right-5 sm:right-8 absolute transform"
				data-testid="clear-button"
			>
				<Icon
					icon="mingcute:close-fill"
					class="text-2xl text-primary"
				/>
			</button>
		</div>

		<div
			class="flex bg-white rounded-full min-w-min overflow-x-auto"
			data-testid="category-buttons"
		>
			<button
				v-for="category in categories"
				:key="category.id"
				:class="[
          'text-base sm:text-xl rounded-full p-5 sm:px-8 sm:py-6 whitespace-nowrap transition-colors font-medium',
          selectedCategoryId === category.id ? 'bg-primary text-primary-foreground' : 'bg-transparent text-gray-500'
        ]"
				@click.prevent="selectCategory(category.id)"
				data-testid="category-button"
			>
				{{ category.name }}
			</button>
		</div>
	</section>
</template>

<script setup lang="ts">
import type { ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { Icon } from '@iconify/vue'
import { ref } from 'vue'


const emit = defineEmits(['update:category', 'update:searchTerm']);
const { selectedCategoryId, searchTerm, categories } = defineProps<{
  selectedCategoryId: number | null;
  searchTerm: string;
  categories: ProductCategoryDTO[];
}>();

const isInputFocused = ref(false);


const updateSearchTerm = (event: Event) => {
  const newTerm = (event.target as HTMLInputElement).value;
  emit('update:searchTerm', newTerm);
};

const clearSearchTerm = () => {
  emit('update:searchTerm', '');
  isInputFocused.value = false;
};

const handleBlur = () => {
  if (!searchTerm.trim()) {
    setTimeout(() => {
      isInputFocused.value = false;
    }, 0);
  }
};

const handleFocus = () => {
  isInputFocused.value = true;
};

const selectCategory = (categoryId: number) => {
  emit('update:category', categoryId);
  emit('update:searchTerm', '');
};
</script>

<style scoped></style>
