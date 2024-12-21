import jsonLocalesMessages from '@/core/locales'
import { createI18n } from 'vue-i18n'

export const LOCALES = {
	RU: 'ru',
	KK: 'kk',
	EN: 'en',
}
export type LocaleTypes = 'ru' | 'kk' | 'en'
export const SUPPORTED_LOCALES = Object.values(LOCALES)

export const AppTranslation = {
	get defaultLocale() {
		return LOCALES.RU
	},

	get supportedLocales() {
		return SUPPORTED_LOCALES
	},

	get currentLocale() {
		return i18nConfig.global.locale.value
	},

	set currentLocale(newLocale) {
		i18nConfig.global.locale.value = newLocale
	},

	async switchLanguage(newLocale: LocaleTypes) {
		AppTranslation.currentLocale = newLocale
		localStorage.setItem('user-locale', newLocale)
	},

	isLocaleSupported(locale: LocaleTypes) {
		return AppTranslation.supportedLocales.includes(locale)
	},

	getUserLocale() {
		const locale = window.navigator.language || AppTranslation.defaultLocale

		return {
			locale: locale,
		}
	},

	getPersistedLocale() {
		const persistedLocale = localStorage.getItem('user-locale') as LocaleTypes

		if (persistedLocale && AppTranslation.isLocaleSupported(persistedLocale)) {
			return persistedLocale
		} else {
			return null
		}
	},

	guessDefaultLocale() {
		const userPersistedLocale = AppTranslation.getPersistedLocale()
		if (userPersistedLocale) {
			return userPersistedLocale
		}

		const userPreferredLocale = AppTranslation.getUserLocale()

		if (AppTranslation.isLocaleSupported(userPreferredLocale.locale as LocaleTypes)) {
			return userPreferredLocale.locale
		}

		return AppTranslation.defaultLocale
	},
}

export const i18nConfig = createI18n({
	locale: AppTranslation.guessDefaultLocale(),
	legacy: false,
	globalInjection: true,
	fallbackLocale: LOCALES.RU,
	messages: jsonLocalesMessages,
})
