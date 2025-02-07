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

		const { beforePrint, afterPrint } = options || {}
		const pdfUrl = URL.createObjectURL(content)

		let cleanupCalled = false
		const cleanup = () => {
			if (cleanupCalled) return
			cleanupCalled = true

			URL.revokeObjectURL(pdfUrl)
			if (afterPrint) afterPrint()
		}

		try {
			if (beforePrint) beforePrint()

			// Only detect Android. If true, open in a new tab.
			if (isAndroid()) {
				/**
				 * Many Android devices/browsers will not show a print dialog
				 * and instead will just download the file when we use `print()`.
				 * We open the PDF in a new tab to at least allow the user to print manually.
				 */
				const newTab = window.open(pdfUrl, '_blank', 'noopener,noreferrer')
				if (!newTab) {
					throw new Error('Failed to open a new tab for printing.')
				}

				// Attempt to call print on load
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
				// Default approach for non-Android devices: print via an iframe
				const iframe = document.createElement('iframe')
				iframe.style.position = 'absolute'
				iframe.style.width = '0'
				iframe.style.height = '0'
				iframe.style.border = '0'
				iframe.style.visibility = 'hidden'
				iframe.src = pdfUrl

				document.body.appendChild(iframe)

				iframe.onload = () => {
					try {
						const iframeWindow = iframe.contentWindow
						if (!iframeWindow) {
							throw new Error('Failed to access iframe window.')
						}
						iframeWindow.focus()
						iframeWindow.print()

						// Cleanup after user closes the print dialog
						iframeWindow.onafterprint = () => {
							cleanup()
							// Remove iframe from DOM
							if (iframe.parentNode) {
								iframe.parentNode.removeChild(iframe)
							}
						}
					} catch (error) {
						console.error('Iframe Print Error:', error)
						if (iframe.parentNode) {
							iframe.parentNode.removeChild(iframe)
						}
						// As a last resort, open the PDF in a new tab
						window.open(pdfUrl, '_blank', 'noopener,noreferrer')
						cleanup()
					}
				}
			}
		} catch (error) {
			console.error('Print Error:', error)
			throw error
		}
	}

	return { print }
}
