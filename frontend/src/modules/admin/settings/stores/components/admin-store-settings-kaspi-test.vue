<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { Switch } from '@/core/components/ui/switch'
import Textarea from '@/core/components/ui/textarea/Textarea.vue'
import { useToast } from '@/core/components/ui/toast'
import { computed, onMounted, onUnmounted, ref } from 'vue'

import {
  KaspiService,
  type KaspiConfig,
  type KaspiDeviceInfo,
  type KaspiTransactionStatus
} from '@/core/integrations/kaspi.service'

// =========================================
// Setup
// =========================================

const { toast } = useToast()

// --- Load Kaspi Config from localStorage
const kaspiConfig = ref<KaspiConfig | null>(null)
onMounted(() => {
  const savedConfig = localStorage.getItem('ZEEP_KASPI_CONFIG')
  if (savedConfig) {
    try {
      kaspiConfig.value = JSON.parse(savedConfig)
    } catch (error) {
      console.warn('Ошибка загрузки конфигурации:', error)
      toast({
        title: '❌ Ошибка загрузки',
        description: 'Не удалось загрузить настройки интеграции.',
        variant: 'destructive'
      })
    }
  }
})

// --- Initialize Kaspi Service if config is available
const kaspi = computed(() => (kaspiConfig.value ? new KaspiService(kaspiConfig.value) : null))

// Helper to clear intervals safely
function clearIntervalIfExists(intervalRef: ReturnType<typeof setInterval> | null) {
  if (intervalRef) {
    clearInterval(intervalRef)
  }
}

// Helper to copy JSON/text to clipboard
const copyToClipboard = (text: string) => {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      toast({
        title: 'Скопировано',
        description: 'JSON успешно скопирован в буфер обмена'
      })
    })
    .catch((error) => {
      toast({
        title: 'Ошибка копирования',
        description: error instanceof Error ? error.message : 'Неизвестная ошибка',
        variant: 'destructive'
      })
    })
}

// =========================================
// Payment Logic
// =========================================

const paymentAmount = ref(100)
const ownCheque = ref<boolean>(true)
const paymentResult = ref<string | null>(null)

// Loading/polling
const paymentLoading = ref(false)
let paymentPolling: ReturnType<typeof setInterval> | null = null

const handlePayment = async () => {
  if (!kaspi.value) return
  if (paymentAmount.value <= 0) {
    toast({
      title: '❌ Ошибка',
      description: 'Сумма оплаты должна быть больше 0.',
      variant: 'destructive'
    })
    return
  }

  try {
    paymentLoading.value = true
    const result = await kaspi.value.initiatePayment({
      amount: paymentAmount.value,
      owncheque: ownCheque.value
    })

    paymentResult.value = `Оплата запущена с Process ID: ${result.data.processId}`
    startPaymentPolling(result.data.processId)

    toast({
      title: '✅ Оплата запущена',
      description: `Process ID: ${result.data.processId}`
    })
  } catch (error) {
    toast({
      title: '❌ Ошибка оплаты',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive'
    })
  } finally {
    paymentLoading.value = false
  }
}

const startPaymentPolling = (processId: string) => {
  clearIntervalIfExists(paymentPolling)

  paymentPolling = setInterval(async () => {
    try {
      const status = await kaspi.value?.getTransactionStatus(processId)
      if (!status) return

      if (status.data.status === 'success') {
        paymentResult.value = JSON.stringify(status, null, 2)
        toast({
          title: '✅ Оплата успешна',
          description: 'Транзакция завершена со статусом success.'
        })
        clearIntervalIfExists(paymentPolling)
      } else if (status.data.status === 'fail') {
        paymentResult.value = JSON.stringify(status, null, 2)
        toast({
          title: '❌ Оплата не удалась',
          description: 'Транзакция завершена со статусом fail.',
          variant: 'destructive'
        })
        clearIntervalIfExists(paymentPolling)
      }
    } catch (error) {
      toast({
        title: '❌ Ошибка опроса',
        description: error instanceof Error ? error.message : 'Неизвестная ошибка',
        variant: 'destructive'
      })
    }
  }, 1000)
}

// =========================================
// Status Check Logic
// =========================================

const statusProcessId = ref('')
const statusResult = ref<KaspiTransactionStatus | null>(null)
const statusResultString = computed(() => JSON.stringify(statusResult))
const statusLoading = ref(false)

const checkStatus = async () => {
  if (!kaspi.value || !statusProcessId.value) return

  try {
    statusLoading.value = true
    const result = await kaspi.value.getTransactionStatus(statusProcessId.value)
    statusResult.value = result

    toast({
      title: '✅ Статус обновлен',
      description: `Статус: ${result.data.status}`
    })
  } catch (error) {
    toast({
      title: '❌ Ошибка статуса',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive'
    })
  } finally {
    statusLoading.value = false
  }
}

const actualizeStatus = async () => {
  if (!kaspi.value || !statusProcessId.value) return

  try {
    statusLoading.value = true
    const result = await kaspi.value.actualizeTransaction(statusProcessId.value)
    statusResult.value = result

    toast({
      title: '✅ Актуализация выполнена',
      description: `Статус: ${result.data.status}`
    })
  } catch (error) {
    toast({
      title: '❌ Ошибка актуализации',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive'
    })
  } finally {
    statusLoading.value = false
  }
}

// =========================================
// Refund Logic
// =========================================

const refundAmount = ref(50)
const refundMethod = ref<'qr' | 'card'>('qr')
const refundTransactionId = ref('')
const refundOwnCheque = ref<boolean>(true)
const refundResult = ref<string | null>(null)

// Loading/polling
const refundLoading = ref(false)
let refundPolling: ReturnType<typeof setInterval> | null = null

const handleRefund = async () => {
  if (!kaspi.value) return

  if (refundAmount.value <= 0) {
    toast({
      title: '❌ Ошибка',
      description: 'Сумма возврата должна быть больше 0.',
      variant: 'destructive'
    })
    return
  }
  if (!refundTransactionId.value) {
    toast({
      title: '❌ Ошибка',
      description: 'Введите Transaction ID.',
      variant: 'destructive'
    })
    return
  }

  try {
    refundLoading.value = true
    const processId = await kaspi.value.refundPayment({
      amount: refundAmount.value,
      method: refundMethod.value,
      transactionId: refundTransactionId.value,
      owncheque: refundOwnCheque.value
    })

    refundResult.value = `Возврат запущен с Process ID: ${processId}`
    startRefundPolling(processId)

    toast({
      title: '✅ Возврат запущен',
      description: `Process ID: ${processId}`
    })
  } catch (error) {
    toast({
      title: '❌ Ошибка возврата',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive'
    })
  } finally {
    refundLoading.value = false
  }
}

const startRefundPolling = (processId: string) => {
  clearIntervalIfExists(refundPolling)

  refundPolling = setInterval(async () => {
    try {
      const status = await kaspi.value?.getTransactionStatus(processId)
      if (!status) return

      if (status.data.status === 'success') {
        refundResult.value = JSON.stringify(status, null, 2)
        toast({
          title: '✅ Возврат успешен',
          description: 'Транзакция завершена со статусом success.'
        })
        clearIntervalIfExists(refundPolling)
      } else if (status.data.status === 'fail') {
        refundResult.value = JSON.stringify(status, null, 2)
        toast({
          title: '❌ Возврат не выполнен',
          description: 'Транзакция завершена со статусом fail.',
          variant: 'destructive'
        })
        clearIntervalIfExists(refundPolling)
      }
    } catch (error) {
      toast({
        title: '❌ Ошибка опроса',
        description: error instanceof Error ? error.message : 'Неизвестная ошибка',
        variant: 'destructive'
      })
    }
  }, 1000)
}

// =========================================
// Device Info Logic
// =========================================

const deviceInfo = ref<KaspiDeviceInfo | null>(null)
const deviceLoading = ref(false)
const deviceResult = ref<string | null>(null)

const getDevice = async () => {
  if (!kaspi.value) return
  try {
    deviceLoading.value = true
    const info = await kaspi.value.getDeviceInfo()
    deviceInfo.value = info
    deviceResult.value = JSON.stringify(info, null, 2)

    toast({
      title: '✅ Информация получена',
      description: `Терминал ID: ${info.data.terminalId}`
    })
  } catch (error) {
    toast({
      title: '❌ Ошибка устройства',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive'
    })
  } finally {
    deviceLoading.value = false
  }
}

// =========================================
// Cleanup intervals
// =========================================

onUnmounted(() => {
  clearIntervalIfExists(paymentPolling)
  clearIntervalIfExists(refundPolling)
})
</script>

<template>
	<div class="space-y-6">
		<!-- Check if we have a valid config -->
		<div
			v-if="!kaspiConfig"
			class="bg-red-50 p-4 border border-red-300 rounded-lg text-red-800 text-center"
		>
			❌ Настройки не найдены. Перейдите в настройки, чтобы указать IP-адрес и имя интеграции.
		</div>

		<template v-else>
			<!-- Payment Section -->
			<div class="space-y-4 bg-gray-50 shadow-sm p-4 border border-gray-200 rounded-lg">
				<h3 class="pb-2 border-b font-semibold text-lg">1. Тест оплаты</h3>
				<div class="flex flex-col space-y-4">
					<div class="space-y-1">
						<Label>Сумма оплаты (₸)</Label>
						<Input
							v-model.number="paymentAmount"
							type="number"
							min="1"
							placeholder="Сумма в тенге"
							class="bg-white max-w-xs"
						/>
					</div>

					<div class="flex items-center space-x-2">
						<Label>Свой чек</Label>
						<Switch
							:model-value="ownCheque"
							@update:model-value="ownCheque = !ownCheque"
						/>
					</div>

					<Button
						:disabled="paymentLoading"
						class="bg-green-600 hover:bg-green-700 w-fit"
						@click="handlePayment"
					>
						<span v-if="!paymentLoading">Начать оплату</span>
						<span v-else>⏳ Выполняется...</span>
					</Button>
				</div>

				<!-- Payment Result TextArea -->
				<div
					v-if="paymentResult"
					class="relative"
				>
					<Textarea
						:value="paymentResult"
						readonly
						class="bg-white pr-14 h-48"
					/>
					<Button
						v-if="paymentResult"
						class="top-2 right-2 absolute"
						size="sm"
						variant="outline"
						@click="copyToClipboard(paymentResult)"
					>
						Копировать JSON
					</Button>
				</div>
			</div>

			<!-- Status Check Section -->
			<div class="space-y-4 bg-gray-50 shadow-sm p-4 border border-gray-200 rounded-lg">
				<h3 class="pb-2 border-b font-semibold text-lg">2. Проверка статуса</h3>
				<div
					class="flex sm:flex-row flex-col items-start sm:items-end sm:space-x-4 space-y-2 sm:space-y-0"
				>
					<div class="space-y-1">
						<Label>Process ID</Label>
						<Input
							v-model="statusProcessId"
							placeholder="Введите Process ID"
							class="bg-white max-w-xs"
						/>
					</div>

					<div class="flex space-x-2">
						<Button
							:disabled="statusLoading || !statusProcessId"
							class="bg-purple-600 hover:bg-purple-700"
							@click="checkStatus"
						>
							Проверить статус
						</Button>

						<Button
							v-if="statusResult?.data.status === 'unknown'"
							:disabled="statusLoading"
							variant="secondary"
							@click="actualizeStatus"
						>
							Актуализировать
						</Button>
					</div>
				</div>

				<!-- Status Result TextArea -->
				<div
					v-if="statusResult"
					class="relative"
				>
					<Textarea
						:value="statusResultString"
						readonly
						class="bg-white pr-14 h-48"
					/>
					<Button
						class="top-2 right-2 absolute"
						size="sm"
						variant="outline"
						@click="copyToClipboard(statusResultString)"
					>
						Копировать JSON
					</Button>
				</div>
			</div>

			<!-- Refund Section -->
			<div class="space-y-4 bg-gray-50 shadow-sm p-4 border border-gray-200 rounded-lg">
				<h3 class="pb-2 border-b font-semibold text-lg">3. Тест возврата</h3>
				<div class="gap-4 grid grid-cols-1">
					<div class="space-y-1">
						<Label>Сумма возврата (₸)</Label>
						<Input
							v-model.number="refundAmount"
							type="number"
							min="1"
							placeholder="50"
							class="bg-white max-w-xs"
						/>
					</div>
					<div class="space-y-1">
						<Label>Метод оплаты</Label>
						<Select v-model="refundMethod">
							<SelectTrigger class="bg-white max-w-xs">
								<SelectValue placeholder="Выберите метод" />
							</SelectTrigger>
							<SelectContent class="bg-white max-w-xs">
								<SelectItem value="qr">QR</SelectItem>
								<SelectItem value="card">Карта</SelectItem>
							</SelectContent>
						</Select>
					</div>

					<div class="flex-1 space-y-1">
						<Label>Transaction ID</Label>
						<Input
							v-model="refundTransactionId"
							placeholder="Введите Transaction ID"
							class="bg-white max-w-xs"
						/>
					</div>

					<div class="flex items-center space-x-2">
						<Label>Свой чек</Label>
						<Switch
							:model-value="refundOwnCheque"
							@update:model-value="refundOwnCheque = !refundOwnCheque"
						/>
					</div>

					<Button
						:disabled="refundLoading"
						class="bg-red-600 hover:bg-red-700 w-fit"
						@click="handleRefund"
					>
						<span v-if="!refundLoading">Начать возврат</span>
						<span v-else>⏳ Выполняется...</span>
					</Button>
				</div>

				<!-- Refund Result TextArea -->
				<div
					v-if="refundResult"
					class="relative"
				>
					<Textarea
						:value="refundResult"
						readonly
						class="bg-white pr-14 h-48"
					/>
					<Button
						class="top-2 right-2 absolute"
						size="sm"
						variant="outline"
						@click="copyToClipboard(refundResult)"
					>
						Копировать JSON
					</Button>
				</div>
			</div>

			<!-- Device Info Section -->
			<div class="space-y-4 bg-gray-50 shadow-sm p-4 border border-gray-200 rounded-lg">
				<h3 class="pb-2 border-b font-semibold text-lg">4. Информация о терминале</h3>
				<div>
					<Button
						:disabled="deviceLoading"
						variant="outline"
						@click="getDevice"
					>
						Получить данные
					</Button>
				</div>

				<!-- Device Info TextArea -->
				<div
					v-if="deviceResult"
					class="relative"
				>
					<Textarea
						:value="deviceResult"
						readonly
						class="bg-white pr-14 h-48"
					/>
					<Button
						class="top-2 right-2 absolute"
						size="sm"
						variant="outline"
						@click="copyToClipboard(deviceResult)"
					>
						Копировать JSON
					</Button>
				</div>
			</div>
		</template>
	</div>
</template>
