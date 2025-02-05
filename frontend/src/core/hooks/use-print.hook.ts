// usePrinter.ts

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
	 * If true, the PDF prints automatically without showing the print dialog.
	 */
	immediatePrint?: boolean
}

export interface UsePrinter {
	/**
	 * Prints the given content (PDF Blob only).
	 * @param content The PDF Blob to print.
	 * @param options Optional printing configurations.
	 */
	print: (content: Blob, options?: PrintOptions) => Promise<void>
}

export function usePrinter(): UsePrinter {
	const print = async (content: Blob, options?: PrintOptions): Promise<void> => {
		if (!(content instanceof Blob) || content.type !== 'application/pdf') {
			throw new Error('Invalid content: Only PDF blobs are supported.')
		}

		const { beforePrint, afterPrint, immediatePrint = false } = options || {}

		try {
			if (beforePrint) beforePrint()

			// Create an object URL for the PDF blob
			const pdfUrl = URL.createObjectURL(content)

			// Create a hidden iframe
			const iframe = document.createElement('iframe')
			iframe.style.position = 'absolute'
			iframe.style.width = '0'
			iframe.style.height = '0'
			iframe.style.border = '0'
			iframe.style.visibility = 'hidden'
			document.body.appendChild(iframe)

			// Load the PDF in the iframe
			iframe.src = pdfUrl

			iframe.onload = () => {
				const iframeWindow = iframe.contentWindow
				if (!iframeWindow) {
					throw new Error('Failed to access iframe window.')
				}

				iframeWindow.focus()

				if (immediatePrint) {
					// 🚀 Attempt seamless printing without a dialog
					iframeWindow.print()
					iframeWindow.onafterprint = () => {
						cleanup()
					}
				} else {
					// Normal print flow
					iframeWindow.print()
					iframeWindow.onafterprint = () => {
						cleanup()
					}
				}
			}

			// Cleanup function
			const cleanup = () => {
				document.body.removeChild(iframe)
				URL.revokeObjectURL(pdfUrl) // Free memory
				if (afterPrint) afterPrint()
			}
		} catch (error) {
			console.error('Print Error:', error)
			throw error
		}
	}

	return { print }
}
