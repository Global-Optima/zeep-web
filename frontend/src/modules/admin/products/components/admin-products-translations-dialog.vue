<template>
	<div>
		<!-- A simple button to open the translation dialog -->
		<Button
			@click="openDialog"
			variant="outline"
			size="icon"
			type="button"
		>
			<Globe
				class="size-4"
				:stroke-width="1.4"
			/>
		</Button>

		<!-- The AdminTranslationDialog component is used here.
         It is bound via v-model for its open state, receives our mock fields,
         and emits a submit event with the DTO payload. -->
		<AdminTranslationsDialog
			v-model:open="isDialogOpen"
			:fields="fields"
			@submit="handleSubmit"
		/>
	</div>
</template>

<script setup lang="ts">
import AdminTranslationsDialog from '@/core/components/admin-translations-dialog/AdminTranslationsDialog.vue'
import { Button } from '@/core/components/ui/button'
import { Globe } from 'lucide-vue-next'
import { ref } from 'vue'

// Define an enum for supported languages.
enum Language {
  EN = 'en',
  KK = 'kk',
  RU = 'ru',
}

// Define the structure for a field with translation information.
interface FieldLocale {
  field: string;
  label: string;
  locales: Partial<Record<Language, string>>;
}

// Some mock data to simulate product translations.
const fields: FieldLocale[] = [
  {
    field: 'name',
    label: 'Название',
    locales: {
      en: 'Product',
      kk: 'Өнім',
      ru: 'Продукт',
    },
  },
  {
    field: 'description',
    label: 'Описание',
    locales: {
      en: 'A great product',
      kk: 'Керемет өнім',
      ru: 'Отличный продукт',
    },
  },
  {
    field: 'address',
    label: 'Адрес',
    locales: {
      en: 'A great product',
      kk: 'Керемет өнім',
      ru: 'Отличный продукт',
    },
  },
  {
    field: 'city',
    label: 'Город',
    locales: {
      en: 'A great product',
      kk: 'Керемет өнім',
      ru: 'Отличный продукт',
    },
  },
  {
    field: 'customer',
    label: 'Клиент',
    locales: {
      en: 'A great product',
      kk: 'Керемет өнім',
      ru: 'Отличный продукт',
    },
  },
]

// Control the visibility of the dialog.
const isDialogOpen = ref(false)

// Open the translation dialog.
function openDialog() {
  isDialogOpen.value = true
}

// Handle the submit event from the dialog.
// The payload is in the desired DTO format, e.g.:
// [ {name: {en: "Product", kk: "Өнім", ru: "Продукт"}}, {description: {...}}, ... ]
function handleSubmit(payload: Record<string, Partial<Record<Language, string>>>[]) {
  console.log('Submitted Translations:', payload)
  // Here you would normally process or send the payload to your API.
  isDialogOpen.value = false
}
</script>

<style scoped>
/* Add any additional component-specific styles here if needed. */
</style>
