// src/composables/useLocalizedToast.ts
import { useToast } from '@/core/components/ui/toast'
import { useI18n } from 'vue-i18n'
import { SUPPORTED_LOCALES } from '../config/locale.config'
import type { LocalizedMessage } from '../models/localized.model'

export type LocalizedToastOptions = {
	variant?: 'default' | 'destructive' | 'success'
	duration?: number
}

type ValidLocale = keyof LocalizedMessage

const isValidLocale = (value: string): value is ValidLocale => {
	return SUPPORTED_LOCALES.includes(value)
}

export const useLocalizedToast = () => {
	const { locale } = useI18n()
	const { toast } = useToast()

	/**
	 * Displays localized toast
	 * @param localizedMessage - Localized message object
	 * @param defaultMessage - Fallback message if no localized message found
	 * @param options - Additional toast configuration
	 */
	const toastLocalized = (
		localizedMessage: LocalizedMessage,
		defaultMessage: string,
		options: LocalizedToastOptions = {},
	) => {
		const currentLocale = locale.value
		const { variant = 'default', duration = 2000 } = options

		const message = isValidLocale(currentLocale) ? localizedMessage[currentLocale] : defaultMessage

		toast({
			variant,
			duration,
			description: message,
		})
	}

	return { toastLocalized }
}
