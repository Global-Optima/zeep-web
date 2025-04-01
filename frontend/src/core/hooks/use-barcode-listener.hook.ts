import { onMounted, onUnmounted } from 'vue'

export interface BarcodeScannerOptions {
	/**
	 * The minimum time (in ms) between characters
	 * to consider them as part of the same scan.
	 * Default: 30 ms.
	 */
	minCharTime?: number

	/**
	 * Barcode prefix (e.g., 'ABC') that must be present
	 * for a successful read. If set and not found, the scan is ignored.
	 */
	prefix?: string

	/**
	 * Barcode suffix (e.g., '\n') that ends a scan.
	 * Default: '\n'
	 */
	suffix?: string

	/**
	 * Callback triggered on successful scan.
	 * @param scannedValue - The final, trimmed barcode.
	 */
	onScan?: (scannedValue: string) => void

	/**
	 * Callback triggered when an error occurs (e.g. parse error).
	 */
	onError?: (error: Error) => void
}

export function useScannerListener(options: BarcodeScannerOptions) {
	const {
		minCharTime = 30, // Adjust if scanner sends too fast/slow
		prefix = '',
		suffix = '\n',
		onScan,
		onError,
	} = options

	// Internal buffer to store scanned characters
	let buffer = ''
	let lastTime = 0
	let timeoutId: number | null = null

	/**
	 * Handles keydown events when scanner inputs barcode.
	 */
	function handleKeyDown(event: KeyboardEvent) {
		// Ignore non-character keys (Shift, Ctrl, Backspace, etc.)
		if (event.key.length !== 1 && event.key !== 'Enter') return

		// Capture the current time
		const currentTime = Date.now()
		const timeDiff = currentTime - lastTime
		lastTime = currentTime

		// If too much time passed, reset the buffer
		if (timeDiff > minCharTime * 2) {
			buffer = ''
		}

		// If Enter is pressed, finalize the barcode
		if (event.key === 'Enter') {
			processBuffer()
			return
		}

		// Append character to buffer
		buffer += event.key

		// Restart timeout to finalize the buffer
		if (timeoutId) {
			clearTimeout(timeoutId)
		}
		timeoutId = window.setTimeout(() => {
			processBuffer()
		}, minCharTime * 5)
	}

	/**
	 * Processes the scanned buffer and triggers `onScan`
	 */
	function processBuffer() {
		if (!buffer.length) return

		try {
			console.debug('Raw buffer:', buffer)

			// Ensure the buffer starts with the required prefix (if specified)
			if (prefix && !buffer.startsWith(prefix)) {
				console.debug('Prefix mismatch:', buffer)
				buffer = ''
				return
			}

			// Remove prefix if present
			let scanned = prefix ? buffer.slice(prefix.length) : buffer

			// Remove suffix if present
			if (suffix && scanned.endsWith(suffix)) {
				scanned = scanned.slice(0, -suffix.length)
			}

			// Ensure we have a valid barcode to process
			if (scanned.length > 0) {
				console.debug('Processed barcode:', scanned)
				onScan?.(scanned)
			}
		} catch (err) {
			if (onError && err instanceof Error) {
				onError(err)
			}
		} finally {
			// Clear the buffer
			buffer = ''
		}
	}

	/**
	 * Lifecycle hooks
	 */
	onMounted(() => {
		window.addEventListener('keydown', handleKeyDown)
	})

	onUnmounted(() => {
		window.removeEventListener('keydown', handleKeyDown)
		if (timeoutId) {
			clearTimeout(timeoutId)
		}
	})
}
