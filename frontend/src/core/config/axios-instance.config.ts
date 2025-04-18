import axios from 'axios'
import { AppTranslation, type LocaleTypes } from './locale.config'

const apiUrl = `${import.meta.env.VITE_API_URL || 'http://example:8080'}/api/v1`

export const apiClient = axios.create({
	baseURL: apiUrl,
	timeout: 60000,
	headers: {
		'Content-Type': 'application/json',
	},
	withCredentials: true,
})

apiClient.interceptors.request.use(
	config => {
		const token = localStorage.getItem('AUTH_TOKEN')
		if (token) {
			config.headers['Authorization'] = `Bearer ${token}`
		}

		const persistedLocale = localStorage.getItem('ZEEP_LOCALE') as LocaleTypes
		if (persistedLocale && AppTranslation.isLocaleSupported(persistedLocale)) {
			config.headers['Accept-Language'] = persistedLocale
		}

		return config
	},
	error => {
		return Promise.reject(error)
	},
)

apiClient.interceptors.response.use(
	response => response,
	async error => {
		if (error.response && error.response.status === 401) {
			console.error('Unauthorized access - redirecting to login')
		}
		return Promise.reject(error)
	},
)
