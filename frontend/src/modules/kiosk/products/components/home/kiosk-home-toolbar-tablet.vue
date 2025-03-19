<template>
	<div
		class="relative"
		data-testid="search-input-container"
	>
		<input
			type="text"
			:value="searchTerm"
			@input="updateSearchTerm"
			:class="['w-full text-base sm:text-xl rounded-3xl p-5 sm:px-8 sm:py-6 focus:outline-none transition-all font-medium']"
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
				class="text-3xl text-primary"
			/>
		</button>
	</div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const emit = defineEmits(['update:category', 'update:searchTerm']);
const { searchTerm } = defineProps<{
  searchTerm: string;
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
</script>

<style scoped></style>
