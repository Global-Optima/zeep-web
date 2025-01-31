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
		/**
		 * 1) Only process keystrokes if the user is focused on
		 *    the element with id="barcode".
		 */
		const activeElem = document.activeElement as HTMLElement | null
		if (!activeElem) {
			// If we are NOT in the barcode input, ignore.
			return
		}

		// 2) Continue with the barcode logic
		const currentTime = Date.now()
		const timeDiff = currentTime - lastTime
		lastTime = currentTime

		// If too much time passed, reset the buffer.
		if (timeDiff > minCharTime) {
			buffer = ''
		}

		// For non-printable keys (like Shift, Ctrl, Arrow keys),
		// we only handle the Enter key to finalize the buffer.
		if (event.key.length !== 1) {
			if (event.key === 'Enter') {
				processBuffer()
			}
			return
		}

		buffer += event.key

		// If we have an existing timeout, clear it and start a new one.
		if (timeout) {
			clearTimeout(timeout)
		}

		// If no further keys after a short delay, treat it as a complete scan.
		timeout = window.setTimeout(() => {
			processBuffer()
		}, minCharTime * 5)
	}

	const processBuffer = () => {
		if (!buffer.length) return

		let barcode = buffer

		// Remove prefix if present
		if (prefix && barcode.startsWith(prefix)) {
			barcode = barcode.substring(prefix.length)
		}
		// Remove suffix if present
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
