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

	// You can adjust these if your device is actually http, or a different port, etc.
	private deviceProtocol = 'https' // was "https://${posIpAddress}:8080"
	private devicePort = '8080'

	constructor(config: KaspiConfig) {
		/**
		 * Instead of calling the device IP directly, we call the local agent:
		 *   http://localhost:42999/proxy
		 *
		 * The local agent will read the query params:
		 *   ?ip=${config.posIpAddress}&port=8080&proto=https
		 *
		 * and forward to: https://192.168.x.x:8080/whatever
		 */
		this.api = axios.create({
			baseURL: 'http://localhost:42999/proxy',
			timeout: config.timeout || 60000,
		})

		// Store the device IP in a property if needed
		this.deviceIp = config.posIpAddress

		this.loadToken()

		// Interceptor to handle 403 => refresh token
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

	private deviceIp: string

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

	/**
	 * Build the query string for local agent:
	 *   ?ip=192.168.200.42&port=8080&proto=https
	 */
	private buildDeviceQuery(params?: Record<string, unknown>): string {
		// Basic device info:
		const baseParams = new URLSearchParams({
			ip: this.deviceIp,
			port: this.devicePort,
			proto: this.deviceProtocol,
		})

		// Add additional query params (like 'name', 'processId', etc.)
		if (params) {
			for (const [key, val] of Object.entries(params)) {
				if (val !== undefined && val !== null) {
					baseParams.append(key, String(val))
				}
			}
		}

		return `?${baseParams.toString()}`
	}

	async registerTerminal(name: string): Promise<KaspiTokenResponse> {
		try {
			// => GET /register?ip=xxx&port=8080&proto=https&name=someName
			const url = `/register${this.buildDeviceQuery({ name })}`
			const response = await this.api.get<KaspiTokenResponse>(url)
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
			// => GET /revoke?ip=xxx&port=8080&proto=https&name=cashier&refreshToken=xxx
			const url = `/revoke${this.buildDeviceQuery({
				name: this.cashierName,
				refreshToken: this.tokenData.data.refreshToken,
			})}`
			const response = await this.api.get<KaspiTokenResponse>(url)
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
			// => GET /payment?ip=xxx&port=8080&proto=https&amount=123&owncheque=true
			const url = `/payment${this.buildDeviceQuery({
				amount: params.amount,
				owncheque: params.owncheque,
			})}`
			const response = await this.api.get<PaymentResponse>(url, {
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка при инициации платежа')
		}
	}

	async getTransactionStatus(processId: string): Promise<KaspiTransactionStatus> {
		try {
			// => GET /status?ip=xxx&port=8080&proto=https&processId=xyz
			const url = `/status${this.buildDeviceQuery({ processId })}`
			const response = await this.api.get<KaspiTransactionStatus>(url, {
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
			// => GET /refund?ip=xxx&port=8080&proto=https&amount=123&method=card&transactionId=abc
			const url = `/refund${this.buildDeviceQuery({
				amount: params.amount,
				method: params.method,
				transactionId: params.transactionId,
				owncheque: params.owncheque,
			})}`
			const response = await this.api.get<RefundResponse>(url, {
				headers: this.getAuthHeaders(),
			})
			return response.data.data.processId
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка возврата платежа')
		}
	}

	async getDeviceInfo(): Promise<KaspiDeviceInfo> {
		try {
			// => GET /deviceinfo?ip=xxx&port=8080&proto=https
			const url = `/deviceinfo${this.buildDeviceQuery()}`
			const response = await this.api.get<KaspiDeviceInfo>(url, {
				headers: this.getAuthHeaders(),
			})
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка получения информации о терминале')
		}
	}

	async actualizeTransaction(processId: string): Promise<KaspiTransactionStatus> {
		try {
			// => GET /actualize?ip=xxx&port=8080&proto=https&processId=abc
			const url = `/actualize${this.buildDeviceQuery({ processId })}`
			const response = await this.api.get<KaspiTransactionStatus>(url, {
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
			`${context ? `${context}: ` : ''}${errorMessage}${
				statusCode ? ` (Код ошибки: ${statusCode})` : ''
			}`,
		)
	}
}
