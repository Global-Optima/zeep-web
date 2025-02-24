import { usePrinter, type PrintOptions } from '@/core/hooks/use-print.hook'
import JsBarcode from 'jsbarcode'
import { jsPDF } from 'jspdf'
import { SAVED_BARISTA_BARCODE_SETTINGS_KEY } from '../constants/barcodes.constants'

export const getSavedBarcodeSettings = (): { width: number; height: number } => {
	const stored = localStorage.getItem(SAVED_BARISTA_BARCODE_SETTINGS_KEY)

	if (stored) {
		try {
			const parsed = JSON.parse(stored)

			if (
				parsed &&
				typeof parsed === 'object' &&
				typeof parsed.width === 'number' &&
				typeof parsed.height === 'number'
			) {
				return { width: parsed.width, height: parsed.height }
			}
		} catch (error) {
			console.error('Error parsing barcode settings from localStorage:', error)
		}
	}

	return { width: 100, height: 50 }
}

export function useGenerateBarcode() {
	/**
	 * Generates a barcode as high-res PNG and embeds it into a PDF (correctly scaled for printing).
	 * @param productName - The name of the product to display above/below the barcode.
	 * @param barcodeValue - The string value to encode in the barcode.
	 * @param labelWidthMm - Width of the label in millimeters.
	 * @param labelHeightMm - Height of the label in millimeters.
	 * @param dpi - Printer DPI (default: 203 DPI for thermal printers).
	 */
	async function generateBarcode(
		productName: string,
		barcodeValue: string,
		labelWidthMm: number = 100,
		labelHeightMm: number = 50,
		dpi: number = 300,
	): Promise<Blob> {
		// -------------------------------
		// 1) Define Label Size in mm
		// -------------------------------
		const pxPerMm = dpi / 25.4 // Conversion factor: DPI to pixels per millimeter

		// Compute usable area inside the label (no hardcoded margins)
		const usableWidthPx = Math.round(labelWidthMm * pxPerMm)
		const usableHeightPx = Math.round(labelHeightMm * pxPerMm)

		// Create a high-res canvas
		const canvas = document.createElement('canvas')
		canvas.width = usableWidthPx
		canvas.height = usableHeightPx

		// -------------------------------
		// 2) Render Barcode (Dynamic Scaling)
		// -------------------------------
		const barcodeHeightPx = Math.round(usableHeightPx * 0.7) // 70% of label height for barcode
		const fontSizePx = Math.round(usableHeightPx * 0.2) // 20% of label height for font size
		const textMarginPx = Math.round(usableHeightPx * 0.02) // 2% of label height for text margin

		JsBarcode(canvas, barcodeValue, {
			format: 'CODE128',
			displayValue: true,
			text: productName,
			font: 'monospace',
			fontSize: fontSizePx,
			textMargin: textMarginPx,
			width: Math.max(1, Math.round(usableWidthPx * 0.01)),
			height: barcodeHeightPx,
			margin: 20,
		})

		// -------------------------------
		// 3) Create PDF (Dynamic Orientation)
		// -------------------------------
		const isLandscape = labelWidthMm > labelHeightMm
		const doc = new jsPDF({
			orientation: isLandscape ? 'landscape' : 'portrait',
			unit: 'mm',
			format: [labelWidthMm, labelHeightMm],
		})

		// Convert the barcode canvas to data URL
		const dataURL = canvas.toDataURL('image/png', 100)

		// -------------------------------
		// 4) Insert Image into PDF (Centered)
		// -------------------------------
		const pdfWidthMm = labelWidthMm
		const pdfHeightMm = labelHeightMm

		// Center the barcode image
		const xStartMm = (pdfWidthMm - labelWidthMm) / 2
		const yStartMm = (pdfHeightMm - labelHeightMm) / 2

		doc.addImage(dataURL, 'PNG', xStartMm, yStartMm, labelWidthMm, labelHeightMm)

		// Return as a Blob for printing
		return doc.output('blob')
	}

	return { generateBarcode }
}

export function useBarcodePrinter() {
	const { print } = usePrinter()
	const { generateBarcode } = useGenerateBarcode()

	/**
	 * Generates and prints a barcode with dynamic label sizes.
	 * @param productName - The name to display with the barcode.
	 * @param barcodeValue - The encoded barcode value.
	 * @param labelWidthMm - Label width in mm.
	 * @param labelHeightMm - Label height in mm.
	 * @param options - Optional print configurations.
	 */
	const printBarcode = async (
		productName: string,
		barcodeValue: string,
		labelWidthMm: number = 100,
		labelHeightMm: number = 50,
		options?: PrintOptions,
	) => {
		try {
			const barcodePdfBlob = await generateBarcode(
				productName,
				barcodeValue,
				labelWidthMm,
				labelHeightMm,
			)
			await print(barcodePdfBlob, options)
		} catch (error) {
			console.error('Error generating or printing barcode:', error)
			throw error
		}
	}

	return { printBarcode }
}
