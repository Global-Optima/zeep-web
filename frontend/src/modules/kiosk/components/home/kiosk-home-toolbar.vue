<template>
	<section
		class="w-full py-4 sm:py-6 px-4 flex items-center gap-2 overflow-x-auto no-scrollbar sticky top-0z-10"
	>
		<div class="relative flex-grow">
			<input
				type="text"
				:value="searchTerm"
				@input="updateSearchTerm"
				@focus="handleFocus"
				@blur="handleBlur"
				:class="[
          'text-base sm:text-xl rounded-full p-5 sm:px-8 sm:py-6 focus:outline-none transition-all font-medium',
          {
            'w-28 sm:w-32 text-center': !isInputFocused,
            'w-80': isInputFocused,
          }
        ]"
				placeholder="Поиск"
			/>

			<button
				v-if="searchTerm"
				@click="clearSearchTerm"
				class="absolute right-5 top-5 sm:right-8 sm:top-6 transform"
			>
				<Icon
					icon="mingcute:close-fill"
					class="text-primary text-2xl"
				/>
			</button>
		</div>

		<div class="flex min-w-min rounded-full overflow-x-auto bg-white">
			<button
				v-for="category in categories"
				:key="category"
				:class="[
          'text-base sm:text-xl rounded-full p-5 sm:px-8 sm:py-6  whitespace-nowrap transition-colors font-medium',
          selectedCategory === category ? 'bg-primary text-primary-foreground' : 'bg-transparent text-gray-500'
        ]"
				@click.prevent="selectCategory(category)"
			>
				{{ category }}
			</button>
		</div>
	</section>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue/dist/iconify.js'
import { ref } from 'vue'

const emit = defineEmits(['update:category', 'update:searchTerm']);
const { selectedCategory, searchTerm } = defineProps<{ selectedCategory: string; searchTerm: string }>();

const categories = [
	'Популярное',
	'Новинки',
	'Сезонное',
	'Кофе',
	'Чай',
	'Холодные напитки',
	'Мороженое'
];

const isInputFocused = ref(false);

const updateSearchTerm = (event: Event) => {
	const newTerm = (event.target as HTMLInputElement).value;
	emit('update:searchTerm', newTerm);
};

const clearSearchTerm = () => {
	emit('update:searchTerm', '');
  isInputFocused.value = false;
}

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

const selectCategory = (category: string) => {
	emit('update:category', category);
	emit('update:searchTerm', '');
};
</script>

<style scoped></style>
