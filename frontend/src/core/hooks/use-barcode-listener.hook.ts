import { onMounted, onUnmounted, ref, type Ref } from 'vue'

interface BarcodeScannerOptions {
	minCharTime?: number
	prefix?: string
	suffix?: string
	onScan: (barcode: string) => Promise<void> | void
	onError?: (error: Error) => void
}

export function useBarcodeScanner(options: BarcodeScannerOptions): {
	lastScanned: Ref<string | null>
} {
	const { minCharTime = 30, prefix = '', suffix = '\n', onScan, onError } = options

	const lastScanned = ref<string | null>(null)
	let buffer = ''
	let lastTime = 0
	let timeout: number | null = null

	const handleKeyDown = (event: KeyboardEvent) => {
		const currentTime = new Date().getTime()
		const timeDiff = currentTime - lastTime
		lastTime = currentTime

		if (timeDiff > minCharTime) {
			buffer = ''
		}

		if (event.key.length !== 1) {
			if (event.key === 'Enter') {
				processBuffer()
			}
			return
		}

		buffer += event.key

		if (timeout) {
			clearTimeout(timeout)
		}

		timeout = window.setTimeout(() => {
			processBuffer()
		}, minCharTime * 5)
	}

	const processBuffer = () => {
		if (buffer.length === 0) return

		let barcode = buffer
		if (prefix && barcode.startsWith(prefix)) {
			barcode = barcode.substring(prefix.length)
		}
		if (suffix && barcode.endsWith(suffix)) {
			barcode = barcode.slice(0, -suffix.length)
		}

		if (barcode.length > 0) {
			try {
				onScan(barcode)
				lastScanned.value = barcode
			} catch (err) {
				if (onError && err instanceof Error) {
					onError(err)
				}
			}
		}

		buffer = ''
	}

	onMounted(() => {
		window.addEventListener('keydown', handleKeyDown)
	})

	onUnmounted(() => {
		window.removeEventListener('keydown', handleKeyDown)
		if (timeout) {
			clearTimeout(timeout)
		}
	})

	return { lastScanned }
}
