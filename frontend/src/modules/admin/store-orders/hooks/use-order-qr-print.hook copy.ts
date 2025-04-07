import { SAVED_BARISTA_QR_SETTINGS_KEY } from '@/core/constants/qr.constants'
import { usePrinter } from '@/core/hooks/use-print.hook'
import type { QRPrintOptions } from '@/core/hooks/use-qr-print.hook'
import type { OrderDTO, SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { MACHINE_CATEGORY_FORMATTED } from '@/modules/kiosk/products/models/product.model'
import { jsPDF } from 'jspdf'
import QRCode from 'qrcode-generator'

// Optionally import vfs_fonts to support Cyrillic characters
// import 'pdfmake/build/vfs_fonts'

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
	// Default fallback (e.g. 80×80 mm label)
	return { width: 80, height: 80 }
}

/**
 * Renders text using an offscreen high‑resolution canvas (which properly supports Cyrillic)
 * and centers each rendered text image within the available width.
 *
 * @param doc jsPDF instance
 * @param text The text to render.
 * @param x Left margin (in mm) for the text block.
 * @param y Vertical starting position (in mm).
 * @param maxWidth Maximum width (in mm) available for text.
 * @param fontSize Font size (in pt).
 * @param isBold Whether to render the text in bold.
 * @returns The updated vertical position after rendering the text.
 */
function drawText(
	doc: jsPDF,
	text: string,
	x: number,
	y: number,
	maxWidth: number,
	fontSize: number,
	isBold: boolean = false,
): number {
	const scaleFactor = 2 // for high-res rendering
	const baseCanvas = document.createElement('canvas')
	const baseCtx = baseCanvas.getContext('2d')!
	if (!baseCtx) return y

	const fontStyle = isBold ? 'bold' : 'normal'
	baseCtx.font = `${fontStyle} ${fontSize}pt Arial`

	// Configure line spacing – increase this factor for more space between lines.
	const lineSpacingFactor = 1
	const lineHeight = fontSize * lineSpacingFactor

	const words = text.split(' ')
	let currentLine = ''
	let currentY = y

	// Process each word to fit within the maxWidth.
	for (let i = 0; i < words.length; i++) {
		const testLine = currentLine + words[i] + ' '
		const metrics = baseCtx.measureText(testLine)
		const testLineWidthMm = metrics.width * 0.35 // approximate conversion px → mm

		if (testLineWidthMm > maxWidth && currentLine !== '') {
			// Render the current line before starting a new one.
			const lineMetrics = baseCtx.measureText(currentLine)
			const textWidthPx = lineMetrics.width

			const textCanvas = document.createElement('canvas')
			textCanvas.width = textWidthPx * scaleFactor
			textCanvas.height = fontSize * 2 * scaleFactor
			const textCtx = textCanvas.getContext('2d')!
			textCtx.scale(scaleFactor, scaleFactor)
			textCtx.font = `${fontStyle} ${fontSize}pt Arial`
			textCtx.fillStyle = '#000000'
			textCtx.textBaseline = 'top'
			textCtx.fillText(currentLine, 0, 0)

			const textDataUrl = textCanvas.toDataURL('image/png')
			const finalTextWidthMm = (textCanvas.width / scaleFactor) * 0.35
			const imageX = x + (maxWidth - finalTextWidthMm) / 2

			doc.addImage(textDataUrl, 'PNG', imageX, currentY, finalTextWidthMm, lineHeight)
			currentY += lineHeight
			currentLine = words[i] + ' '
		} else {
			currentLine = testLine
		}
	}

	// Render any remaining text.
	if (currentLine.trim()) {
		const lineMetrics = baseCtx.measureText(currentLine)
		const textWidthPx = lineMetrics.width

		const textCanvas = document.createElement('canvas')
		textCanvas.width = textWidthPx * scaleFactor
		textCanvas.height = fontSize * 2 * scaleFactor
		const textCtx = textCanvas.getContext('2d')!
		textCtx.scale(scaleFactor, scaleFactor)
		textCtx.font = `${fontStyle} ${fontSize}pt Arial`
		textCtx.fillStyle = '#000000'
		textCtx.textBaseline = 'top'
		textCtx.fillText(currentLine, 0, 0)

		const textDataUrl = textCanvas.toDataURL('image/png')
		const finalTextWidthMm = (textCanvas.width / scaleFactor) * 0.35
		const imageX = x + (maxWidth - finalTextWidthMm) / 2

		doc.addImage(textDataUrl, 'PNG', imageX, currentY, finalTextWidthMm, lineHeight)
		currentY += lineHeight
	}
	return currentY
}

export function useOrderGenerateQR() {
	/**
	 * Generates a QR code as a high-res PNG and embeds it into a PDF (scaled for printing),
	 * along with two lines of text:
	 *   1) "#OrderNumber  CustomerName"
	 *   2) "ProductName ProductSizeString (CategoryName, 350 ml)"
	 */
	async function generateQRFromSuborder(
		order: OrderDTO,
		suborder: SuborderDTO,
		options: QRPrintOptions = {},
	): Promise<Blob> {
		// 1) Retrieve label dimensions and DPI (defaulting if necessary)
		const { dpi = 203, labelWidthMm = 80, labelHeightMm = 80 } = options
		const pxPerMm = dpi / 25.4

		// 2) Prepare the QR content.
		const qrValue = generateSubOrderQR(suborder)

		// 3) Define text lines.
		const line1 = `#${order.displayNumber} ${order.customerName}`
		const line2 = `${suborder.productSize.productName} ${suborder.productSize.sizeName} (${suborder.productSize.size} ${suborder.productSize.unit.name}, ${MACHINE_CATEGORY_FORMATTED[suborder.productSize.machineCategory]})`

		// 4) Generate the QR code onto a canvas.
		const usableWidthPx = Math.round(labelWidthMm * pxPerMm)
		const usableHeightPx = Math.round(labelHeightMm * pxPerMm)

		const canvas = document.createElement('canvas')
		canvas.width = usableWidthPx
		canvas.height = usableHeightPx

		const qr = QRCode(0, 'M')
		qr.addData(qrValue)
		qr.make()

		const ctx = canvas.getContext('2d')!
		ctx.fillStyle = '#FFFFFF'
		ctx.fillRect(0, 0, canvas.width, canvas.height)
		ctx.fillStyle = '#000000'

		const qrModules = qr.getModuleCount()
		const qrScaleFactor = 0.9 // increased scale for a larger, high-res QR
		const qrSizePx = Math.min(usableWidthPx, usableHeightPx) * qrScaleFactor
		const cellSize = qrSizePx / qrModules
		const qrMarginPx = (Math.min(usableWidthPx, usableHeightPx) - qrSizePx) / 2

		for (let row = 0; row < qrModules; row++) {
			for (let col = 0; col < qrModules; col++) {
				if (qr.isDark(row, col)) {
					ctx.fillRect(qrMarginPx + col * cellSize, qrMarginPx + row * cellSize, cellSize, cellSize)
				}
			}
		}

		// Convert the QR canvas to a PNG data URL.
		const qrDataURL = canvas.toDataURL('image/png', 1.0)

		// 5) Create a new jsPDF document with the proper dimensions.
		const doc = new jsPDF({
			orientation: labelWidthMm > labelHeightMm ? 'landscape' : 'portrait',
			unit: 'mm',
			format: [labelWidthMm, labelHeightMm],
		})

		// 6) Determine a responsive font size.
		const baseFontSize = Math.min(10, labelWidthMm / 8)

		// 7) Define margins.
		const margin = Math.max(2, labelWidthMm * 0.025)
		let currentY = margin

		// 8) Calculate the maximum text width (accounting for margins).
		const maxTextWidth = labelWidthMm - margin * 2

		// 9) Render the text lines using our high-res canvas approach.
		currentY = drawText(doc, line1, margin, currentY, maxTextWidth, baseFontSize, true)
		currentY = drawText(doc, line2, margin, currentY, maxTextWidth, baseFontSize * 0.9)

		// 10) Add additional spacing between the text and QR code.
		const qrSpacing = baseFontSize * 0.5
		currentY += qrSpacing

		// 11) Determine the QR code’s size and position.
		const remainingHeight = labelHeightMm - currentY - margin
		const qrImageSizeMm = Math.min(remainingHeight, labelWidthMm - 2 * margin)
		const finalQrSizeMm = Math.max(qrImageSizeMm, Math.min(labelWidthMm, labelHeightMm) * 0.3)
		const qrX = (labelWidthMm - finalQrSizeMm) / 2

		if (currentY + finalQrSizeMm <= labelHeightMm - margin) {
			doc.addImage(qrDataURL, 'PNG', qrX, currentY, finalQrSizeMm, finalQrSizeMm)
		} else {
			const adjustedSize = labelHeightMm - currentY - margin
			if (adjustedSize > 10) {
				doc.addImage(qrDataURL, 'PNG', qrX, currentY, adjustedSize, adjustedSize)
			}
		}

		// 12) Return the PDF as a Blob.
		return doc.output('blob')
	}

	return { generateQRFromSuborder }
}

export function useOrderQRPrinter() {
	const { print } = usePrinter()
	const { generateQRFromSuborder } = useOrderGenerateQR()

	/**
	 * Prints one or more QR labels for a suborder.
	 * Accepts either a single SuborderDTO or an array.
	 */
	const printSubOrderQR = async (order: OrderDTO, options?: QRPrintOptions) => {
		try {
			const suborders = Array.isArray(order.subOrders) ? order.subOrders : [order.subOrders]
			const pdfBlobs = await Promise.all(
				suborders.map(sb => generateQRFromSuborder(order, sb, options)),
			)
			await print(pdfBlobs, options)
		} catch (error) {
			console.error('Error generating or printing QR code:', error)
			throw error
		}
	}

	return { printSubOrderQR }
}
