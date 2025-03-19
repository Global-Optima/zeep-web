<template>
	<Select v-model="selectedLocale">
		<SelectTrigger
			class="bg-gray-600/60 backdrop-blur-md px-6 py-8 border-none rounded-3xl text-base text-white sm:text-3xl"
		>
			<SelectValue
				placeholder="Select language"
				class="flex items-center gap-3"
			/>
		</SelectTrigger>
		<SelectContent>
			<SelectItem
				v-for="locale in supportedLocales"
				:key="locale.locale"
				:value="locale.locale"
			>
				<div class="flex items-center gap-4">
					<Icon
						:icon="locale.icon"
						class="text-2xl sm:text-3xl"
					/>
					<span class="text-base sm:text-3xl">{{ locale.label }}</span>
				</div>
			</SelectItem>
		</SelectContent>
	</Select>
</template>

<script setup lang="ts">
import { AppTranslation, LOCALES, type LocaleTypes } from '@/core/config/locale.config'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'



// Import Iconify for icons
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import { Icon } from '@iconify/vue'

const { locale } = useI18n()
const router = useRouter()

const supportedLocales = [
  {
    label: 'EN',
    locale: LOCALES.EN,
    icon: 'twemoji:flag-united-kingdom',
  },
  {
    label: 'RU',
    locale: LOCALES.RU,
    icon: 'twemoji:flag-russia',
  },
  {
    label: 'KZ',
    locale: LOCALES.KK,
    icon: 'twemoji:flag-kazakhstan',
  },
]

const selectedLocale = ref(locale.value)

const switchLanguage = async (newLocale: LocaleTypes) => {
  await AppTranslation.switchLanguage(newLocale)
  try {
    await router.replace({ params: { locale: newLocale } })
  } catch (e) {
    console.error('Error switching language:', e)
    await router.push('/')
  }
}

// Watch for changes in selectedLocale and switch language
watch(selectedLocale, (newLocale) => {
  locale.value = newLocale
  switchLanguage(newLocale as LocaleTypes)
})
</script>

<style scoped></style>
