import { SAVED_BARISTA_QR_SETTINGS_KEY } from '@/core/constants/qr.constants'
import type { SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { jsPDF } from 'jspdf'
import QRCode from 'qrcode-generator'
import { usePrinter, type PrintOptions } from './use-print.hook'

export interface QRPrintOptions extends PrintOptions {
	labelWidthMm?: number
	labelHeightMm?: number
	dpi?: number
}

export const generateSubOrderQR = (sb: SuborderDTO) => {
	const parts = [
		sb.id ?? null,
		sb.productSize?.machineId ?? null,
		sb.additives?.length ? sb.additives.map(a => a.additive.machineId).join(',') : null,
	]

	return parts.filter(Boolean).join('|')
}

export const parseSubOrderQR = (qrString: string) => {
	const [subOrderId, productSizeMachineId, additivesString] = qrString.split('|')

	return {
		subOrderId: subOrderId || null,
		productSizeMachineId: productSizeMachineId || null,
		additiveMachineIds: additivesString ? additivesString.split(',') : [],
	}
}

export const getSavedBaristaQRSettings = (): { width: number; height: number } => {
	const stored = localStorage.getItem(SAVED_BARISTA_QR_SETTINGS_KEY)

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
			console.error('Error parsing QR settings from localStorage:', error)
		}
	}

	return { width: 100, height: 100 } // Default to a square QR label
}

export function useGenerateQR() {
	/**
	 * Generates a QR code as high-res PNG and embeds it into a PDF (correctly scaled for printing).
	 * @param productName - The name of the product to display above/below the QR code.
	 * @param qrValue - The string value to encode in the QR code.
	 * @param labelWidthMm - Width of the label in millimeters.
	 * @param labelHeightMm - Height of the label in millimeters.
	 * @param dpi - Printer DPI (default: 300 for high-res prints).
	 */
	async function generateQR(qrValue: string, options: QRPrintOptions = {}): Promise<Blob> {
		const { dpi = 203, labelWidthMm = 100, labelHeightMm = 100 } = options
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
		// 2) Generate QR Code
		// -------------------------------
		const qr = QRCode(0, 'M') // M = Medium Error Correction
		qr.addData(qrValue)
		qr.make()

		// Get QR code cell size
		const qrSizePx = Math.min(usableWidthPx, usableHeightPx) * 0.8 // Scale QR to 80% of label size
		const qrMarginPx = (Math.min(usableWidthPx, usableHeightPx) - qrSizePx) / 2 // Centering

		// Draw QR on canvas
		const ctx = canvas.getContext('2d')!
		ctx.fillStyle = '#FFFFFF' // White background
		ctx.fillRect(0, 0, canvas.width, canvas.height) // Clear canvas
		ctx.fillStyle = '#000000' // Black for QR Code

		const qrModules = qr.getModuleCount()
		const cellSize = qrSizePx / qrModules

		for (let row = 0; row < qrModules; row++) {
			for (let col = 0; col < qrModules; col++) {
				if (qr.isDark(row, col)) {
					ctx.fillRect(qrMarginPx + col * cellSize, qrMarginPx + row * cellSize, cellSize, cellSize)
				}
			}
		}

		// -------------------------------
		// 3) Create PDF (Dynamic Orientation)
		// -------------------------------
		const isLandscape = labelWidthMm > labelHeightMm
		const doc = new jsPDF({
			orientation: isLandscape ? 'landscape' : 'portrait',
			unit: 'mm',
			format: [labelWidthMm, labelHeightMm],
		})

		// Convert the QR canvas to data URL
		const dataURL = canvas.toDataURL('image/png', 100)

		// -------------------------------
		// 4) Insert Image into PDF (Centered)
		// -------------------------------
		const pdfWidthMm = labelWidthMm
		const pdfHeightMm = labelHeightMm

		// Center the QR code image
		const xStartMm = (pdfWidthMm - labelWidthMm) / 2
		const yStartMm = (pdfHeightMm - labelHeightMm) / 2

		doc.addImage(dataURL, 'PNG', xStartMm, yStartMm, labelWidthMm, labelHeightMm)

		// Return as a Blob for printing
		return doc.output('blob')
	}

	return { generateQR }
}

export function useQRPrinter() {
	const { print } = usePrinter()
	const { generateQR } = useGenerateQR()

	/**
	 * @param qrValue - The encoded QR value.
	 * @param labelWidthMm - Label width in mm.
	 * @param labelHeightMm - Label height in mm.
	 * @param options - Optional print configurations.
	 */
	const printQR = async (qrValue: string | string[], options?: QRPrintOptions) => {
		try {
			const qrValues = Array.isArray(qrValue) ? qrValue : [qrValue]
			const qrPdfBlobs = await Promise.all(qrValues.map(v => generateQR(v, options)))

			await print(qrPdfBlobs, options)
		} catch (error) {
			console.error('Error generating or printing QR code:', error)
			throw error
		}
	}

	return { printQR }
}
