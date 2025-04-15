import { saveAs } from 'file-saver'
import isMobile from 'is-mobile'
import printJS from 'print-js'
import { v4 as uuidv4 } from 'uuid'

/**
 * If your .env has VITE_SAVE_ON_PRINT set to 'true',
 * we’ll save files instead of using the print dialog.
 */
const SAVE_ON_PRINT = import.meta.env.VITE_SAVE_ON_PRINT === 'true'

export interface PrintOptions {
	/** If true, printing is allowed only on desktop devices.
	 *  On mobile/tablet, the print job will be skipped.
	 */
	desktopOnly?: boolean
	/** Invoked before printing each file. */
	beforePrint?: () => void
	/** Invoked after printing each file. */
	afterPrint?: () => void

	/** Whether to force 'save' on the file instead of printing. */
	saveOnPrint?: boolean
}

interface PrintJob {
	blobs: Blob[]
	options?: PrintOptions
	resolve: () => void
	reject: (err: unknown) => void
}

export function usePrinter() {
	// Queue of print jobs
	const jobQueue: PrintJob[] = []
	let isProcessing = false

	/**
	 * Main API: queue up a print (or save) job.
	 * The job can contain one or more Blob(s).
	 */
	const print = async (content: Blob | Blob[], options?: PrintOptions): Promise<void> => {
		// 1) If desktopOnly & device is mobile, skip
		if (options?.desktopOnly && isMobile({ tablet: true })) {
			console.warn('[usePrinter] Skipped: Desktop-only mode on a mobile device.')
			return Promise.resolve()
		}

		// 2) Convert content → array of Blobs
		const blobs = Array.isArray(content) ? content : [content]

		// 3) Create a job & enqueue
		return new Promise<void>((resolve, reject) => {
			jobQueue.push({ blobs, options, resolve, reject })
			processQueue()
		})
	}

	/**
	 * Processes jobs in FIFO order, one at a time.
	 */
	async function processQueue() {
		if (isProcessing) return
		isProcessing = true

		while (jobQueue.length > 0) {
			const job = jobQueue.shift()!
			try {
				// Print each Blob in the job sequentially
				for (const blob of job.blobs) {
					await printOrSaveSingleFile(blob, job.options)
				}
				// Mark the entire job as done
				job.resolve()
			} catch (error) {
				job.reject(error)
			}
		}

		isProcessing = false
	}

	/**
	 * Prints or saves a single file, depending on:
	 * - The 'saveOnPrint' option (either local or fallback to global)
	 * - The Blob's type (PDF vs. ZPL or anything else)
	 */
	async function printOrSaveSingleFile(blob: Blob, opts?: PrintOptions) {
		return new Promise<void>((resolve, reject) => {
			const { beforePrint, afterPrint, saveOnPrint = SAVE_ON_PRINT } = opts || {}
			if (typeof beforePrint === 'function') beforePrint()

			// Create an object URL for printing
			const fileUrl = URL.createObjectURL(blob)
			let cleanedUp = false

			// Cleanup helper
			const cleanup = () => {
				if (!cleanedUp) {
					cleanedUp = true
					URL.revokeObjectURL(fileUrl)
					if (typeof afterPrint === 'function') afterPrint()
				}
			}

			try {
				if (saveOnPrint) {
					const uniqueFileName = `zeep-qr-${uuidv4()}.prn`
					console.warn(
						`[usePrinter] Non-PDF mime type detected: ${blob.type}. Saving file instead of printing.`,
					)
					saveAs(blob, uniqueFileName)
					cleanup()
					resolve()
				} else {
					// Option 2: We want to "print", but we must handle PDF vs. non‑PDF
					if (blob.type === 'application/pdf') {
						// Print PDF with printJS
						printJS({
							printable: fileUrl,
							type: 'pdf',
							onLoadingEnd: cleanup,
							onPrintDialogClose: () => {
								console.log('[usePrinter] PDF print dialog closed.')
								resolve()
							},
							// If you want to debug or pass header info, you can add more config here
						})
					}
				}
			} catch (err) {
				cleanup()
				reject(err)
			}
		})
	}

	return { print }
}
