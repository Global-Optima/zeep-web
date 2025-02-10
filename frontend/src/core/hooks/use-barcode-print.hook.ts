import { usePrinter, type PrintOptions } from '@/core/hooks/use-print.hook'
import JsBarcode from 'jsbarcode'
import { jsPDF } from 'jspdf'

export function useGenerateBarcode() {
	/**
	 * Generates a barcode as high-res PNG and embeds into PDF (landscape).
	 * With added margins so it fits nicely and doesn’t get squeezed.
	 */
	async function generateBarcode(productName: string, barcodeValue: string): Promise<Blob> {
		// ----------------------------------
		// 1) Define physical size in mm
		//    3.88 cm × 2.6 cm => 38.8 mm × 26.0 mm
		//    orientation: 'landscape' => jsPDF format is [height, width]
		// ----------------------------------
		const labelWidthMm = 38.8
		const labelHeightMm = 26.0

		// We’ll add a margin (e.g., 2 mm) inside the label.
		// This keeps the barcode from touching edges or being cut off.
		const marginMm = 4

		// Calculate the “printable” area inside the margins.
		const usableWidthMm = labelWidthMm - marginMm * 2
		const usableHeightMm = labelHeightMm - marginMm * 2

		// ----------------------------------
		// 2) Create a high-res canvas to generate the barcode
		// ----------------------------------
		// ~3.78 px per mm is a rough conversion.
		// Multiply further by a “scaleFactor” for crispness.
		const pxPerMm = 3.78
		const scaleFactor = 3 // you can try 2, 3, 4, etc.

		// Convert the usable mm to pixels (minus margins).
		const canvasWidthPx = Math.round(usableWidthMm * pxPerMm * scaleFactor)
		const canvasHeightPx = Math.round(usableHeightMm * pxPerMm * scaleFactor)

		const canvas = document.createElement('canvas')
		canvas.width = canvasWidthPx
		canvas.height = canvasHeightPx

		// ----------------------------------
		// 3) Render the barcode smaller so it fits
		//    Decrease width, height, and textMargin if needed
		// ----------------------------------
		JsBarcode(canvas, barcodeValue, {
			format: 'CODE128',
			displayValue: true,
			text: productName,
			font: 'monospace',
			fontSize: 14 * scaleFactor,
			width: 1.5 * scaleFactor,
			height: Math.round(canvasHeightPx * 0.5),
			margin: 4 * scaleFactor,
		})

		// ----------------------------------
		// 4) Create jsPDF instance with the correct size in mm, landscape
		// ----------------------------------
		const doc = new jsPDF({
			orientation: 'landscape',
			unit: 'mm',
			format: [labelHeightMm, labelWidthMm],
		})

		// Convert the barcode canvas to data URL
		const dataURL = canvas.toDataURL('image/png')

		// ----------------------------------
		// 5) Insert the barcode image with margin into the PDF
		//    So it doesn’t fill the entire label edge-to-edge
		// ----------------------------------
		doc.addImage(
			dataURL,
			'PNG',
			marginMm, // x start
			marginMm, // y start
			usableWidthMm, // fit into the “usable” space
			usableHeightMm,
		)

		// Output PDF as Blob
		return doc.output('blob')
	}

	return { generateBarcode }
}

export function useBarcodePrinter() {
	const { print } = usePrinter()
	const { generateBarcode } = useGenerateBarcode()

	/**
	 * Generates and prints a PDF barcode.
	 * @param productName - The name of the product to display above/below the barcode.
	 * @param barcodeValue - The string value to encode in the barcode.
	 * @param options - Optional printing configurations.
	 */
	const printBarcode = async (
		productName: string,
		barcodeValue: string,
		options?: PrintOptions,
	) => {
		try {
			// Just call the generateBarcode hook without passing custom sizes;
			// The function already uses the correct label size + 300 DPI.
			const barcodePdfBlob = await generateBarcode(productName, barcodeValue)
			await print(barcodePdfBlob, options)
		} catch (error) {
			console.error('Error generating or printing barcode:', error)
			throw error
		}
	}

	return { printBarcode }
}
