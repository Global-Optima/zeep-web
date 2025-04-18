<template>
	<Dialog
		:open="isOpen"
		@update:open="onUpdateOpenDialog"
	>
		<DialogContent
			:include-close-button="false"
			class="max-h-[95vh] overflow-y-auto"
		>
			<DialogHeader>
				<DialogTitle> Редактор переводов </DialogTitle>
				<DialogDescription> Управление переводами для ваших полей </DialogDescription>
			</DialogHeader>

			<div
				v-if="loading"
				class="flex justify-center items-center mt-8"
			>
				<Loader class="size-8 text-primary animate-spin" />
			</div>

			<div v-else>
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
										v-for="lang in getSortedLocales(item)"
										:key="lang"
										:class="cn(getLanguageBadgeClass(lang), 'text-xs px-2.5 py-0.5 rounded-md')"
									>
										{{ lang.toUpperCase() }}
									</span>
								</div>
							</div>
						</AccordionTrigger>
						<AccordionContent>
							<div class="space-y-3">
								<!-- Существующие языки, отсортированные по порядку RU, KK, EN -->
								<div
									v-for="lang in getSortedLocales(item)"
									:key="lang"
									class="flex items-center gap-3"
								>
									<span
										:class="cn(getLanguageBadgeClass(lang), 'text-xs px-2.5 py-0.5 rounded-md')"
									>
										{{ lang.toUpperCase() }}
									</span>
									<Input
										v-model="localFields[index].locales[lang]"
										:placeholder="`Введите перевод на ${getLanguageName(lang)}`"
										class="w-full text-sm"
									/>
									<button
										@click="removeLocale(index, lang)"
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
					:disabled="loading"
					>Отмена</Button
				>
				<Button
					@click="submit"
					class="flex items-center"
					:disabled="loading"
				>
					Сохранить
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import { TranslationsLanguage, type TranslationFieldLocale } from '@/core/components/admin-translations-dialog'
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
import Input from '@/core/components/ui/input/Input.vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/core/components/ui/select'
import { cn } from '@/core/utils/tailwind.utils'
import { Loader, PlusIcon, XIcon } from 'lucide-vue-next'
import { reactive, ref, watch } from 'vue'



// Массив языков в требуемом порядке: RU, KK, EN.
const sortedLanguages: TranslationsLanguage[] = [TranslationsLanguage.RU, TranslationsLanguage.KK, TranslationsLanguage.EN]

// Отображаемые названия языков
const languageNames: Record<TranslationsLanguage, string> = {
  [TranslationsLanguage.EN]: 'Английский',
  [TranslationsLanguage.KK]: 'Казахский',
  [TranslationsLanguage.RU]: 'Русский',
}

// Новые цвета для бейджей, ассоциированные с каждым языком
const languageColors: Record<TranslationsLanguage, string> = {
  [TranslationsLanguage.RU]: 'bg-rose-100 text-rose-700',
  [TranslationsLanguage.KK]: 'bg-yellow-100 text-yellow-700',
  [TranslationsLanguage.EN]: 'bg-sky-100 text-sky-700',
}

// Тип для описания поля с переводами.

// Пропсы компонента.
const props = defineProps<{
  fields: TranslationFieldLocale[];
  open: boolean;
  loading: boolean;
}>()

// События компонента.
const emit = defineEmits<{
  (e: 'submit', payload: Record<string, Partial<Record<TranslationsLanguage, string>>>[]): void;
  (e: 'update:open', value: boolean): void;
}>()

// Локальная реактивная копия полей для редактирования.
const localFields = reactive<TranslationFieldLocale[]>(props.fields.map(field => ({
  ...field,
  locales: { ...field.locales },
})))

// Отслеживаем выбор новой локали для каждого поля.
const newLocaleForField = ref<Record<number, TranslationsLanguage | ''>>({})

// Синхронизируем входящие пропсы с локальным состоянием.
watch(() => props.fields, (newFields) => {
  localFields.splice(0, localFields.length, ...newFields.map(field => ({
    ...field,
    locales: { ...field.locales },
  })))
  newLocaleForField.value = {}
}, { deep: true })

// Управление состоянием диалога.
const isOpen = ref(props.open)
watch(() => props.open, (newVal) => {
  isOpen.value = newVal
})

// Возвращает отсортированные локали для поля в порядке RU, KK, EN.
function getSortedLocales(item: TranslationFieldLocale): TranslationsLanguage[] {
  return sortedLanguages.filter(lang => lang in item.locales)
}

// Возвращает недостающие локали для конкретного поля.
function getMissingLocalesForField(index: number): TranslationsLanguage[] {
  const fieldLocales = new Set(Object.keys(localFields[index].locales))
  return sortedLanguages.filter(lang => !fieldLocales.has(lang))
}

// Добавление локали к конкретному полю.
function addLocaleToField(index: number): void {
  const locale = newLocaleForField.value[index]
  if (locale && !localFields[index].locales[locale]) {
    localFields[index].locales[locale] = ''
    // Сброс выбора.
    newLocaleForField.value[index] = ''
  }
}

// Удаление локали у конкретного поля.
function removeLocale(index: number, locale: TranslationsLanguage): void {
  if (localFields[index].locales[locale] !== undefined) {
    delete localFields[index].locales[locale]
  }
}

// Функция для получения отображаемого названия языка.
function getLanguageName(code: TranslationsLanguage): string {
  return languageNames[code] || code
}

// Функция для получения класса бейджа для языка.
function getLanguageBadgeClass(code: TranslationsLanguage): string {
  return languageColors[code] || 'bg-gray-100 text-gray-800'
}

// Обработка обновления состояния диалога.
const onUpdateOpenDialog = (value: boolean): void => {
  emit('update:open', value)
}

// Закрытие диалога.
function closeDialog(): void {
  onUpdateOpenDialog(false)
}

// Подготовка и отправка данных.
function submit(): void {
  const payload = localFields.map(field => ({
    [field.field]: { ...field.locales },
  }))
  emit('submit', payload)
  closeDialog()
}
</script>
