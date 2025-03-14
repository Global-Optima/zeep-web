import type { TransactionDTO } from '@/modules/admin/store-orders/models/orders.models'
import axios, { type AxiosInstance } from 'axios'

export interface KaspiConfig {
	host: string
	name: string
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
		subStatus:
			| 'Initialize'
			| 'WaitUser'
			| 'WaitForQrConfirmation'
			| 'ProcessingCard'
			| 'WaitForPinCode'
			| 'ProcessRefund'
			| 'QrTransactionSuccess'
			| 'QrTransactionFailure'
			| 'CardTransactionSuccess'
			| 'CardTransactionFailure'
			| 'ProcessCancelled'
		transactionId: string
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

export const KASPI_CONFIG_STORAGE_KEY = 'ZEEP_KASPI_CONFIG'
export const KASPI_TOKENS_STORAGE_KEY = 'ZEEP_KASPI_TOKENS'

export const getKaspiConfig = () => {
	const savedConfig = localStorage.getItem(KASPI_CONFIG_STORAGE_KEY)
	if (savedConfig) {
		try {
			const parsedConfig: KaspiConfig = JSON.parse(savedConfig)
			return parsedConfig
		} catch (error) {
			console.warn('Ошибка загрузки конфигурации:', error)
		}
	}
}

export const saveKaspiConfig = (config: KaspiConfig) => {
	localStorage.setItem(KASPI_CONFIG_STORAGE_KEY, JSON.stringify(config))
}

export class KaspiService {
	private api: AxiosInstance
	private tokenData: KaspiTokenResponse | null = null
	private config: KaspiConfig

	constructor(config: KaspiConfig) {
		this.config = config
		this.api = axios.create({
			baseURL: `https://${config.host}:8080/v2`,
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

	async awaitPaymentTest(): Promise<TransactionDTO> {
		await new Promise(resolve => setTimeout(resolve, 5000)) // Simulate delay

		// Randomize success/failure (50% chance each)
		const isSuccess = Math.random() > 0.5

		if (!isSuccess) {
			throw new Error('Mock Payment Failure')
		}

		return {
			bin: '123456',
			transactionId: 'TX-' + Date.now(),
			processId: 'PROC-' + (Math.random() * 1e5).toFixed(0),
			paymentMethod: 'card',
			amount: 1500,
			currency: 'KZT',
			qrNumber: 'QR-' + (Math.random() * 1e6).toFixed(0),
			cardMask: '************4242',
			icc: 'fakeICCData',
		}
	}

	public async awaitPayment(amount: number): Promise<TransactionDTO> {
		try {
			const initiateResponse = await this.initiatePayment({ amount })
			return this.pollPaymentStatus(initiateResponse.data.processId)
		} catch (error) {
			throw new Error(
				`Payment initiation failed: ${error instanceof Error ? error.message : 'Unexpected error'}`,
			)
		}
	}

	private async pollPaymentStatus(processId: string): Promise<TransactionDTO> {
		return new Promise<TransactionDTO>(async (resolve, reject) => {
			const poll = async () => {
				try {
					const statusResponse = await this.getTransactionStatus(processId)
					const { status, message, transactionId, chequeInfo, subStatus } = statusResponse.data

					console.log('STATUSSSS', status)

					if (status === 'success') {
						console.log('SUCCESFULLY START')
						if (!chequeInfo) {
							console.log('SUCCESFULLY NO CHECK INFO')
							reject(new Error('Invalid payment response: Missing cheque information'))
							return
						}

						resolve({
							bin: chequeInfo.bin,
							transactionId: transactionId,
							processId,
							paymentMethod: chequeInfo.method,
							amount: Number(chequeInfo.amount),
							currency: 'KZT',
							qrNumber: chequeInfo.method === 'qr' ? chequeInfo.orderNumber : undefined,
							cardMask: chequeInfo.method === 'card' ? chequeInfo.cardMask : undefined,
							icc: chequeInfo.method === 'card' ? chequeInfo.icc : undefined,
						})

						console.log('SUCCESFULLY AFTER CHECK INFO')

						return
					}

					if (
						status === 'fail' ||
						status === 'unknown' ||
						subStatus === 'QrTransactionFailure' ||
						subStatus === 'CardTransactionFailure' ||
						subStatus === 'ProcessCancelled'
					) {
						reject(new Error(`Payment failed: ${message ?? 'Unexpected error'}`))
						return
					}

					setTimeout(poll, 1000)
				} catch (error) {
					reject(
						new Error(
							`Status check failed: ${error instanceof Error ? error.message : 'Unexpected error'}`,
						),
					)
				}
			}

			poll()
		})
	}

	private saveToken(): void {
		if (this.tokenData) {
			localStorage.setItem(KASPI_TOKENS_STORAGE_KEY, JSON.stringify(this.tokenData))
		}
	}

	private loadToken(): void {
		const savedToken = localStorage.getItem(KASPI_TOKENS_STORAGE_KEY)
		if (savedToken) {
			this.tokenData = JSON.parse(savedToken)
		}
	}

	async registerTerminal(): Promise<KaspiTokenResponse> {
		try {
			const response = await this.api.get<KaspiTokenResponse>('/register', {
				params: { name: this.config.name },
			})

			this.tokenData = response.data
			this.saveToken()
			return response.data
		} catch (error) {
			throw this.handleApiError(error, 'Ошибка регистрации кассового аппарата')
		}
	}

	async refreshAccessToken(): Promise<void> {
		if (!this.tokenData || !this.config.name) {
			throw new Error('Нет действительного токена или имени кассира')
		}

		try {
			const response = await this.api.get<KaspiTokenResponse>('/revoke', {
				params: {
					name: this.config.name,
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
