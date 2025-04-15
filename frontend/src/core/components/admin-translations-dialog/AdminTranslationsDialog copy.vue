<template>
	<Dialog
		:open="isOpen"
		@update:open="onUpdateOpenDialog"
		class="p-4"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle> Редактор переводов </DialogTitle>
				<DialogDescription> Управление переводами для ваших полей </DialogDescription>
			</DialogHeader>

			<!-- Переработанный интерфейс с аккордеоном для каждого поля -->
			<div>
				<Accordion
					type="multiple"
					class="w-full"
				>
					<AccordionItem
						v-for="(item, index) in localFields"
						:key="item.field"
						:value="item.field"
					>
						<AccordionTrigger class="py-6 hover:no-underline">
							<div class="flex justify-between items-center w-full">
								<span class="font-semibold">{{ item.label }}</span>
								<div class="flex gap-1 mr-2">
									<span
										v-for="lang in Object.keys(item.locales)"
										:key="lang"
										:class="cn(getLanguageBadgeClass(lang as Language), 'text-xs px-2.5 py-0.5 rounded-md')"
									>
										{{ lang.toUpperCase() }}
									</span>
								</div>
							</div>
						</AccordionTrigger>
						<AccordionContent>
							<div class="space-y-3">
								<!-- Существующие языки -->
								<div
									v-for="lang in Object.keys(item.locales)"
									:key="lang"
									class="flex items-center gap-3"
								>
									<span
										:class="cn(getLanguageBadgeClass(lang as Language), 'text-xs px-2.5 py-0.5 rounded-md')"
									>
										{{ lang.toUpperCase() }}
									</span>
									<input
										v-model="localFields[index].locales[lang as Language]"
										type="text"
										:placeholder="`Введите перевод на ${getLanguageName(lang as Language)}`"
										class="p-2 border border-gray-300 focus:border-blue-500 rounded focus:ring-blue-500 w-full text-sm"
									/>
									<button
										@click="removeLocale(index, lang as Language)"
										class="hover:bg-gray-100 p-1 rounded-full text-gray-400 hover:text-red-600"
										title="Удалить перевод"
									>
										<XIcon class="w-4 h-4" />
									</button>
								</div>

								<!-- Добавление нового языка для текущего поля -->
								<div
									v-if="getMissingLocalesForField(index).length > 0"
									class="flex justify-center items-center gap-2 mt-2 p-[1px]"
								>
									<Select v-model="newLocaleForField[index]">
										<SelectTrigger class="w-fit">
											<SelectValue placeholder="Выберите язык" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="lang in getMissingLocalesForField(index)"
												:key="lang"
												:value="lang"
											>
												{{ getLanguageName(lang) }}
											</SelectItem>
										</SelectContent>
									</Select>
									<Button
										size="icon"
										variant="outline"
										@click="addLocaleToField(index)"
										:disabled="!newLocaleForField[index]"
									>
										<PlusIcon class="w-4 h-4" />
									</Button>
								</div>
							</div>
						</AccordionContent>
					</AccordionItem>
				</Accordion>
			</div>

			<!-- Кнопки действий -->
			<div class="flex justify-end gap-2 mt-6">
				<Button
					variant="outline"
					@click="closeDialog"
				>
					Отмена
				</Button>
				<Button
					@click="submit"
					class="flex items-center"
				>
					Сохранить
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/core/components/ui/accordion'
import Button from '@/core/components/ui/button/Button.vue'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from '@/core/components/ui/dialog'
import DialogDescription from '@/core/components/ui/dialog/DialogDescription.vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { cn } from '@/core/utils/tailwind.utils'
import { PlusIcon, XIcon } from 'lucide-vue-next'
import { defineEmits, defineProps, reactive, ref, watch } from 'vue'


// Define an enum for the supported languages.
enum Language {
  EN = 'en',
  KK = 'kk',
  RU = 'ru',
}

// List of allowed languages (reusable everywhere)
const allowedLanguages = Object.values(Language);

// Language display names
const languageNames: Record<Language, string> = {
  [Language.EN]: 'Английский',
  [Language.KK]: 'Казахский',
  [Language.RU]: 'Русский',
};

// Language badge colors
const languageColors: Record<Language, string> = {
  [Language.EN]: 'bg-blue-100 text-blue-800',
  [Language.KK]: 'bg-green-100 text-green-800',
  [Language.RU]: 'bg-purple-100 text-purple-800',
};

// Define the shape of a field with locales.
interface FieldLocale {
  field: string;
  label: string;
  locales: {
    [key in Language]?: string;
  };
}

// Define component props.
const props = defineProps<{
  fields: FieldLocale[];
  open: boolean;
}>();

// Define component events.
const emit = defineEmits<{
  (e: 'submit', payload: Record<string, Partial<Record<Language, string>>>[]): void;
  (e: 'update:open', value: boolean): void;
}>();

// Create a reactive local copy of the fields for user editing.
const localFields = reactive<FieldLocale[]>(props.fields.map(field => ({
  ...field,
  locales: { ...field.locales },
})));

// Track new locale selection for each field
const newLocaleForField = ref<Record<number, Language | ''>>({});

// Watch for changes to the incoming props to keep the local copy in sync.
watch(() => props.fields, (newFields) => {
  localFields.splice(0, localFields.length, ...newFields.map(field => ({
    ...field,
    locales: { ...field.locales },
  })));
  newLocaleForField.value = {};
}, { deep: true });

// Control dialog open/close state.
const isOpen = ref(props.open);
watch(() => props.open, (newVal) => {
  isOpen.value = newVal;
});

// Get missing locales for a specific field
function getMissingLocalesForField(index: number): Language[] {
  const fieldLocales = new Set(Object.keys(localFields[index].locales));
  return allowedLanguages.filter(lang => !fieldLocales.has(lang));
}

// Add a locale to a specific field
function addLocaleToField(index: number): void {
  const locale = newLocaleForField.value[index];
  if (locale && !localFields[index].locales[locale]) {
    localFields[index].locales[locale] = '';
    // Clear the selection
    newLocaleForField.value[index] = '';
  }
}

// Remove a locale from a specific field
function removeLocale(index: number, locale: Language): void {
  if (localFields[index].locales[locale] !== undefined) {
    delete localFields[index].locales[locale];
  }
}

// Get language display name
function getLanguageName(code: Language): string {
  return languageNames[code] || code;
}

// Get badge class for language
function getLanguageBadgeClass(code: Language): string {
  return languageColors[code] || 'bg-gray-100 text-gray-800';
}

// Handle dialog update
const onUpdateOpenDialog = (value: boolean): void => {
  emit('update:open', value);
}

// Close the dialog
function closeDialog(): void {
  onUpdateOpenDialog(false);
}

// On submission, transform the local fields into the required DTO format
function submit(): void {
  const payload = localFields.map(field => ({
    [field.field]: { ...field.locales },
  }));
  emit('submit', payload);
  closeDialog();
}
</script>
