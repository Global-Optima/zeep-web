// src/composables/useAxiosLocaleToast.ts
import { AxiosError } from 'axios'
import type { LocalizedError } from '../models/errors.model'
import { useLocalizedToast, type LocalizedToastOptions } from './use-locale-toast.hooks'

export type AxiosLocalizedError = AxiosError<LocalizedError>

const initialDefaultMessage = 'Возникла непредвиденная ошибка. Попробуйте позже.'

export const useAxiosLocaleToast = () => {
	const { toastLocalized } = useLocalizedToast()

	/**
	 * Displays localized error toast based on Axios error
	 * @param error - The Axios error object
	 * @param defaultMessage - Fallback message if no localized message found
	 * @param options - Additional toast configuration
	 */
	const toastLocalizedError = (
		error: AxiosLocalizedError,
		defaultMessage: string,
		options?: LocalizedToastOptions,
	) => {
		if (error.response) {
			const localizedMessage = error.response.data.message
			if (localizedMessage) {
				return toastLocalized(localizedMessage, defaultMessage, {
					...options,
					variant: 'destructive',
				})
			}
		}

		toastLocalized(
			{
				en: 'Network error occurred',
				ru: 'Ошибка сети',
				kk: 'Желі қатесі',
			},
			initialDefaultMessage,
			{ variant: 'destructive' },
		)
	}

	return { toastLocalizedError }
}
