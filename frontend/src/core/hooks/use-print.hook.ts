// usePrinter.ts

export type PrintContent = string | Blob

export interface PrintOptions {
	/**
	 * Optional CSS styles to apply to the printed content.
	 */
	styles?: string

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
	 * Prints the given content (image URL or file Blob).
	 * @param content The image URL or file Blob to print.
	 * @param options Optional printing configurations.
	 */
	print: (content: PrintContent, options?: PrintOptions) => Promise<void>
}

export function usePrinter(): UsePrinter {
	const blobToDataURL = (blob: Blob): Promise<string> => {
		return new Promise((resolve, reject) => {
			const reader = new FileReader()
			reader.onload = () => {
				if (typeof reader.result === 'string') {
					resolve(reader.result)
				} else {
					reject(new Error('Failed to convert blob to Data URL.'))
				}
			}
			reader.onerror = () => {
				reject(new Error('Error reading the blob.'))
			}
			reader.readAsDataURL(blob)
		})
	}

	const isImageURL = (url: string): boolean => {
		return /\.(jpeg|jpg|gif|png|bmp|webp|svg)$/.test(url.toLowerCase())
	}

	const print = async (content: PrintContent, options?: PrintOptions): Promise<void> => {
		const { styles = '', beforePrint, afterPrint } = options || {}

		try {
			if (beforePrint) beforePrint()

			let contentToPrint = ''

			if (typeof content === 'string') {
				if (isImageURL(content)) {
					contentToPrint = `<img src="${content}" style="max-width: 100%; height: auto;" />`
				} else {
					// For non-image URLs, attempt to embed in iframe (e.g., PDF)
					contentToPrint = `<iframe src="${content}" style="width: 100%; height: 100%;" frameborder="0"></iframe>`
				}
			} else if (content instanceof Blob) {
				const dataURL = await blobToDataURL(content)
				const mimeType = content.type
				if (mimeType.startsWith('image/')) {
					contentToPrint = `<img src="${dataURL}" style="max-width: 100%; height: auto;" />`
				} else if (mimeType === 'application/pdf') {
					contentToPrint = `<iframe src="${dataURL}" style="width: 100%; height: 100%;" frameborder="0"></iframe>`
				} else {
					throw new Error(`Unsupported file type: ${mimeType}`)
				}
			} else {
				throw new Error('Unsupported content type.')
			}

			// Create a hidden iframe for printing
			const iframe = document.createElement('iframe')
			iframe.style.position = 'fixed'
			iframe.style.right = '0'
			iframe.style.bottom = '0'
			iframe.style.width = '0'
			iframe.style.height = '0'
			iframe.style.border = '0'
			iframe.style.visibility = 'hidden'
			iframe.id = `print-iframe-${Date.now()}`
			document.body.appendChild(iframe)

			const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document
			if (!iframeDoc) {
				throw new Error('Failed to access iframe document.')
			}

			iframeDoc.open()
			iframeDoc.write(`
        <html>
          <head>
            <title>Print Preview</title>
            <style>
              /* Append any custom styles passed in by the user */
              ${styles}

              /* Ensure no extra page margins */
              @page {
                margin: 0;
              }

              @media print {
                html, body {
                  margin: 0;
                  padding: 0;
                }
                img, iframe {
                  max-width: 100%;
                  height: auto;
                }
              }
            </style>
          </head>
          <body>
            ${contentToPrint}
          </body>
        </html>
      `)
			iframeDoc.close()

			// Wait for the iframe content to load
			await new Promise<void>((resolve, reject) => {
				const iframeWindow = iframe.contentWindow
				if (!iframeWindow) {
					reject(new Error('Failed to access iframe window.'))
					return
				}

				const handleLoad = () => {
					// Trigger the print dialog
					iframeWindow.focus()
					iframeWindow.print()

					// Cleanup after print
					iframeWindow.onafterprint = () => {
						resolve()
					}
				}

				// If the content is an iframe (like a PDF), use .addEventListener('load')
				iframe.addEventListener('load', handleLoad)

				// Fallback if 'load' doesn't fire
				setTimeout(() => {
					resolve()
				}, 5000)
			})

			document.body.removeChild(iframe)

			if (afterPrint) afterPrint()
		} catch (error) {
			console.error('Print Error:', error)
			throw error
		}
	}

	return { print }
}
