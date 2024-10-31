import { nextTick } from 'vue'
import { createI18n, type I18n } from 'vue-i18n'
import type {
	NavigationGuardNext,
	RouteLocationNormalizedGeneric,
	RouteLocationNormalizedLoadedGeneric,
} from 'vue-router'
import ru from '../locales/ru.json'

export const LOCALES = {
	RU: 'ru',
	KK: 'kk',
	EN: 'en',
}
export type LocaleTypes = 'ru' | 'kk' | 'en'
export const SUPPORTED_LOCALES = Object.values(LOCALES)

type I18nType = I18n<
	{
		ru: object
		en: object
		kk: object
	},
	Record<string, unknown>,
	Record<string, unknown>,
	LocaleTypes,
	false
>

export const i18nConfig: I18nType = createI18n({
	locale: LOCALES.RU,
	legacy: false,
	globalInjection: true,
	fallbackLocale: LOCALES.RU,
	messages: { ru },
})

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
		await AppTranslation.loadLocaleMessages(newLocale)
		AppTranslation.currentLocale = newLocale
		document.querySelector('html')?.setAttribute('lang', newLocale)
		localStorage.setItem('user-locale', newLocale)
	},

	async loadLocaleMessages(locale: LocaleTypes) {
		if (!i18nConfig.global.availableLocales.includes(locale)) {
			const messages = await import(`@/core/locales/${locale}.json`)
			i18nConfig.global.setLocaleMessage(locale, messages.default)
		}

		return nextTick()
	},

	isLocaleSupported(locale: LocaleTypes) {
		return AppTranslation.supportedLocales.includes(locale)
	},

	getUserLocale() {
		const locale = window.navigator.language || AppTranslation.defaultLocale

		return {
			locale: locale,
			localeNoRegion: locale.split('-')[0],
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

		if (AppTranslation.isLocaleSupported(userPreferredLocale.localeNoRegion as LocaleTypes)) {
			return userPreferredLocale.localeNoRegion
		}

		return AppTranslation.defaultLocale
	},

	async routeMiddleware(
		to: RouteLocationNormalizedGeneric,
		_from: RouteLocationNormalizedLoadedGeneric,
		next: NavigationGuardNext,
	) {
		const paramLocale = to.params.locale as LocaleTypes

		if (!AppTranslation.isLocaleSupported(paramLocale)) {
			return next(AppTranslation.guessDefaultLocale())
		}

		await AppTranslation.switchLanguage(paramLocale)

		return next()
	},

	i18nRoute(to: RouteLocationNormalizedGeneric) {
		return {
			...to,
			params: {
				locale: AppTranslation.currentLocale,
				...to.params,
			},
		}
	},
}
