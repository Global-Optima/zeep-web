import { saveAs } from 'file-saver'
import isMobile from 'is-mobile'
import printJS from 'print-js'
import { v4 as uuidv4 } from 'uuid'

const SAVE_ON_PRINT = import.meta.env.VITE_SAVE_ON_PRINT === 'true'

export interface PrintOptions {
	/** If true, printing is allowed only on desktop devices.
	 *  On mobile/tablet, the print job will be skipped.
	 */
	desktopOnly?: boolean
	/** Invoked before printing each PDF. */
	beforePrint?: () => void
	/** Invoked after printing each PDF. */
	afterPrint?: () => void

	saveOnPrint?: boolean
}

interface PrintJob {
	blobs: Blob[] // one or more PDF Blobs
	options?: PrintOptions
	resolve: () => void // For resolving the .print(...) Promise
	reject: (err: unknown) => void
}

/**
 * Check if the argument is a valid PDF Blob (application/pdf).
 */
function isPDFBlob(content: unknown): content is Blob {
	return content instanceof Blob && content.type === 'application/pdf'
}

/**
 * Hook for printing PDF blobs in a queued, sequential manner.
 */
export function usePrinter() {
	// Queue of print jobs
	const jobQueue: PrintJob[] = []
	let isProcessing = false

	/**
	 * Enqueues a print job and returns a Promise that resolves
	 * when the entire job (all blobs) has finished printing.
	 */
	const print = async (content: Blob | Blob[], options?: PrintOptions): Promise<void> => {
		// If desktopOnly flag is set and the device is mobile, skip printing.
		if (options?.desktopOnly && isMobile({ tablet: true })) {
			console.warn('[usePrinter] Printing skipped: Desktop-only mode active on a mobile device.')
			return Promise.resolve()
		}

		// Normalize to an array of Blobs
		const blobs = Array.isArray(content) ? content : [content]

		// Validate each blob
		for (const blob of blobs) {
			if (!isPDFBlob(blob)) {
				throw new Error('Invalid content: Only PDF blobs are supported.')
			}
		}

		// Enqueue the job
		return new Promise<void>((resolve, reject) => {
			jobQueue.push({ blobs, options, resolve, reject })
			processQueue()
		})
	}

	/**
	 * Core logic: process the queue in FIFO order.
	 * Only one job is processed at a time. Once finished,
	 * move on to the next job until the queue is empty.
	 */
	async function processQueue() {
		if (isProcessing) return
		isProcessing = true

		while (jobQueue.length > 0) {
			const job = jobQueue.shift()!
			try {
				// Print each Blob in the job sequentially
				for (const blob of job.blobs) {
					await printSinglePDF(blob, job.options)
				}
				// All blobs in this job have been printed
				job.resolve()
			} catch (error) {
				job.reject(error)
			}
		}

		isProcessing = false
	}

	/**
	 * Prints or Saves a single PDF file based on the given options.
	 *
	 * @param blob - The PDF file blob.
	 * @param opts - Printing options.
	 * @returns A promise that resolves once printing or saving is done.
	 */
	async function printSinglePDF(blob: Blob, opts?: PrintOptions): Promise<void> {
		return new Promise<void>((resolve, reject) => {
			const { beforePrint, afterPrint, saveOnPrint = SAVE_ON_PRINT } = opts || {}

			if (typeof beforePrint === 'function') beforePrint()

			let cleanedUp = false
			const pdfUrl = URL.createObjectURL(blob)

			// Cleanup function to revoke the blob URL and call `afterPrint` once
			const cleanup = () => {
				if (!cleanedUp) {
					cleanedUp = true
					URL.revokeObjectURL(pdfUrl)
					if (typeof afterPrint === 'function') afterPrint()
				}
			}

			try {
				// If saveOnPrint is true, skip printing and just trigger a save
				if (saveOnPrint) {
					const uniqueFileName = `${uuidv4()}.pdf`

					saveAs(blob, uniqueFileName)
					cleanup()
					resolve()
				} else {
					// Otherwise, proceed with printing using printJS
					printJS({
						printable: pdfUrl,
						type: 'pdf',
						onLoadingEnd: cleanup, // called once PDF has loaded in print preview
						onPrintDialogClose: () => {
							console.log('[printSinglePDF] Print dialog closed')
							resolve()
						},
					})
				}
			} catch (err) {
				cleanup()
				reject(err)
			}
		})
	}

	return { print }
}
