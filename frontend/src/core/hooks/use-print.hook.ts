import printJS from 'print-js'

export interface PrintOptions {
	/**
	 * Optional callback invoked before printing starts.
	 */
	beforePrint?: () => void

	/**
	 * Optional callback invoked after printing completes.
	 */
	afterPrint?: () => void

	/**
	 * Determines whether Print.js shows its "printing" modal dialog.
	 * Default is `true` if not specified.
	 */
	showModal?: boolean
}

export interface UsePrinter {
	/**
	 * Prints the given content (PDF Blob only).
	 * @param content The PDF Blob to print.
	 * @param options Optional printing configurations.
	 */
	print: (content: Blob, options?: PrintOptions) => Promise<void>
}

function isPDFBlob(content: Blob): boolean {
	return content instanceof Blob && content.type === 'application/pdf'
}

function isAndroid(): boolean {
	return /Android/i.test(navigator.userAgent)
}

export function usePrinter(): UsePrinter {
	const print = async (content: Blob, options?: PrintOptions): Promise<void> => {
		if (!isPDFBlob(content)) {
			throw new Error('Invalid content: Only PDF blobs are supported.')
		}

		const { beforePrint, afterPrint, showModal = true } = options || {}
		const pdfUrl = URL.createObjectURL(content)

		let cleanupCalled = false
		const cleanup = () => {
			if (cleanupCalled) return
			cleanupCalled = true

			URL.revokeObjectURL(pdfUrl)
			if (afterPrint) {
				afterPrint()
			}
		}

		try {
			// 1) Call the beforePrint callback if provided
			if (beforePrint) beforePrint()

			// 2) Handle Android separately (per original logic)
			if (isAndroid()) {
				const newTab = window.open(pdfUrl, '_blank', 'noopener,noreferrer')
				if (!newTab) {
					throw new Error('Failed to open a new tab for printing (Android).')
				}

				// Attempt to print once the tab loads
				newTab.onload = () => {
					try {
						newTab.focus()
						newTab.print()
					} catch (error) {
						console.error('Print Error on new tab:', error)
					} finally {
						cleanup()
					}
				}

				// Fallback if onload doesn’t fire within 2 seconds
				setTimeout(() => {
					if (!cleanupCalled) {
						try {
							newTab.focus()
							newTab.print()
						} catch (error) {
							console.error('Print Error in fallback:', error)
						} finally {
							cleanup()
						}
					}
				}, 2000)
			} else {
				// 3) For non-Android, use Print.js
				printJS({
					printable: pdfUrl,
					type: 'pdf',
					showModal: showModal,
					onLoadingEnd: () => {
						cleanup()
					},
				})
			}
		} catch (error) {
			console.error('Print Error:', error)
			cleanup()
			throw error
		}
	}

	return { print }
}
