import axios, { type AxiosInstance } from 'axios'

export interface KaspiConfig {
	posIpAddress: string
	integrationName: string
	timeout?: number
}

export interface KaspiTokenResponse {
	data: {
		accessToken: string
		refreshToken: string
		expirationDate: string
	}
}

interface PaymentRequest {
	amount: number
	owncheque?: boolean
}

interface PaymentResponse {
	data: {
		processId: string
		status: string
	}
	statusCode: number
}

interface RefundRequest {
	amount: number
	method: 'qr' | 'card'
	transactionId: string
	owncheque?: boolean
}

interface RefundResponse {
	data: {
		processId: string
		status: string
	}
}

export interface KaspiTransactionStatus {
	data: {
		processId: string
		status: 'wait' | 'success' | 'fail' | 'unknown'
		subStatus?: string
		transactionId?: string
		message?: string
		addInfo?: {
			IsOffer?: boolean
			ProductType?: string
			LoanTerm?: number
			LoanOfferName?: string
		}
		chequeInfo?: QRPaymentInfo | CardPaymentInfo
	}
}

interface QRPaymentInfo {
	storeName: string
	city: string
	address: string
	status: string
	amount: string
	date: string
	bin: string
	terminalId: string
	orderNumber: string
	method: 'qr'
}

interface CardPaymentInfo {
	storeName: string
	city: string
	address: string
	status: string
	amount: string
	date: string
	cardMask: string
	icc: string
	bin: string
	terminalId: string
	rrn: string
	authorizationCode: string
	hostResponseCode: string
	method: 'card'
}

export interface KaspiDeviceInfo {
	data: {
		posNum: string
		serialNum: string
		terminalId: string
	}
}

export class KaspiService {
	private api: AxiosInstance
	private tokenData: KaspiTokenResponse | null = null
	private cashierName: string | null = null

	constructor(config: KaspiConfig) {
		this.api = axios.create({
			baseURL: `${window.location.origin}/external/kaspi/?ip=${config.posIpAddress}&port=8080`,
			timeout: config.timeout || 60000,
		})

		this.loadToken()

		this.api.interceptors.response.use(
			response => response,
			async error => {
				const originalRequest = error.config

				if (error.response?.status === 403 && !originalRequest._retry && this.tokenData) {
					originalRequest._retry = true

					try {
						await this.refreshAccessToken()
						return this.api(originalRequest)
					} catch (refreshError) {
						throw new Error(`Ошибка аутентификации: ${refreshError}`)
					}
				}

				throw this.handleApiError(error)
			},
		)
	}

	private saveToken(): void {
		if (this.tokenData) {
			localStorage.setItem('kaspi_token', JSON.stringify(this.tokenData))
		}
	}

	private loadToken(): void {
		const savedToken = localStorage.getItem('kaspi_token')
		if (savedToken) {
			this.tokenData = JSON.parse(savedToken)
		}
	}

	async registerTerminal(name: string): Promise<KaspiTokenResponse> {
		try {
			const response = await this.api.get<KaspiTokenResponse>('/register', { params: { name } })

			this.tokenData = response.data
			this.cashierName = name
			this.saveToken()
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка регистрации кассового аппарата')
		}
	}

	async refreshAccessToken(): Promise<void> {
		if (!this.tokenData || !this.cashierName) {
			throw new Error('Нет действительного токена или имени кассира')
		}

		try {
			const response = await this.api.get<KaspiTokenResponse>('/revoke', {
				params: {
					name: this.cashierName,
					refreshToken: this.tokenData.data.refreshToken,
				},
			})

			this.tokenData = response.data
			this.saveToken()
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка обновления токена')
		}
	}

	async initiatePayment(params: PaymentRequest) {
		if (params.amount <= 0) {
			throw new Error('Сумма платежа должна быть больше нуля')
		}

		try {
			const response = await this.api.get<PaymentResponse>('/payment', {
				params: { amount: params.amount, owncheque: params.owncheque },
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка при инициации платежа')
		}
	}

	async getTransactionStatus(processId: string): Promise<KaspiTransactionStatus> {
		try {
			const response = await this.api.get<KaspiTransactionStatus>('/status', {
				params: { processId },
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка проверки статуса транзакции')
		}
	}

	async refundPayment(params: RefundRequest): Promise<string> {
		if (params.amount <= 0) {
			throw new Error('Сумма возврата должна быть больше нуля')
		}

		try {
			const response = await this.api.get<RefundResponse>('/refund', {
				params: {
					amount: params.amount,
					method: params.method,
					transactionId: params.transactionId,
					owncheque: params.owncheque,
				},
				headers: this.getAuthHeaders(),
			})
			return response.data.data.processId
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка возврата платежа')
		}
	}

	async getDeviceInfo(): Promise<KaspiDeviceInfo> {
		try {
			const response = await this.api.get<KaspiDeviceInfo>('/deviceinfo', {
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка получения информации о терминале')
		}
	}

	// ✅ New Actualization Method
	async actualizeTransaction(processId: string): Promise<KaspiTransactionStatus> {
		try {
			const response = await this.api.get<KaspiTransactionStatus>('/actualize', {
				params: { processId },
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка актуализации транзакции')
		}
	}

	private getAuthHeaders(): Record<string, string> {
		if (!this.tokenData) {
			throw new Error('Необходима аутентификация')
		}
		return { accesstoken: this.tokenData.data.accessToken }
	}

	private handleApiError(error: unknown, context?: string): Error {
		let errorMessage = 'Произошла неизвестная ошибка'
		let statusCode: number | undefined

		if (error instanceof Object && 'response' in error) {
			const axiosError = error as {
				response?: { data?: { errorText?: string; statusCode?: number } }
			}
			errorMessage = axiosError.response?.data?.errorText || 'Ответ сервера отсутствует'
			statusCode = axiosError.response?.data?.statusCode
		} else if (error instanceof Error) {
			errorMessage = error.message
		}

		return new Error(
			`${context ? `${context}: ` : ''}${errorMessage}${statusCode ? ` (Код ошибки: ${statusCode})` : ''}`,
		)
	}
}
