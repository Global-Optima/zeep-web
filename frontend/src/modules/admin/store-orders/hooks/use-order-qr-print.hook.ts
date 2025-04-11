import { SAVED_BARISTA_QR_SETTINGS_KEY } from '@/core/constants/qr.constants'
import { usePrinter } from '@/core/hooks/use-print.hook'
import type { QRPrintOptions } from '@/core/hooks/use-qr-print.hook'
import type { OrderDTO, SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { MACHINE_CATEGORY_FORMATTED } from '@/modules/kiosk/products/models/product.model'
import { jsPDF } from 'jspdf'
import QRCode from 'qrcode-generator'

// If using VITE_SAVE_ON_PRINT in your .env file:
const SAVE_ON_PRINT = import.meta.env.VITE_SAVE_ON_PRINT === 'true'

/**
 * Utility to combine the suborder IDs, product machine IDs, etc. into a single string
 * that you can parse out later (the "qrValue").
 */
export function generateSubOrderQR(sb: SuborderDTO): string {
	if (!sb.id || !sb.productSize?.machineId) {
		throw new Error('suborderId or productSizeMachineId is missing')
	}

	const parts = [
		sb.id,
		sb.productSize.machineId,
		sb.additives?.length ? sb.additives.map(a => a.additive.machineId).join(',') : null,
	]

	return parts.filter(Boolean).join('|')
}

export function parseSubOrderQR(qrString: string) {
	const [subOrderId, productSizeMachineId, additivesString] = qrString.split('|')

	if (!subOrderId || !productSizeMachineId) {
		throw new Error('suborderId or productSizeMachineId is missing')
	}

	return {
		subOrderId: subOrderId,
		productSizeMachineId: productSizeMachineId,
		additiveMachineIds: additivesString ? additivesString.split(',') : [],
	}
}

/**
 * Gets the user's saved label dimensions from localStorage,
 * falling back to 80x80 mm if none is saved.
 */
export function getSavedBaristaQRSettings(): { width: number; height: number } {
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
	return { width: 80, height: 80 }
}

/**
 * This helper draws text onto an offscreen canvas (so it can handle Cyrillic),
 * then places the rasterized text in a PDF at the correct position.
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
	// Text → image logic
	const scaleFactor = 2
	const baseCanvas = document.createElement('canvas')
	const baseCtx = baseCanvas.getContext('2d')!
	if (!baseCtx) return y

	const fontStyle = isBold ? 'bold' : 'normal'
	baseCtx.font = `${fontStyle} ${fontSize}pt Helvetica`

	const lineSpacingFactor = 1
	const lineHeight = fontSize * lineSpacingFactor

	const words = text.split(' ')
	let currentLine = ''
	let currentY = y

	for (let i = 0; i < words.length; i++) {
		const testLine = currentLine + words[i] + ' '
		const metrics = baseCtx.measureText(testLine)
		const testLineWidthMm = metrics.width * 0.35 // approximate px → mm

		if (testLineWidthMm > maxWidth && currentLine !== '') {
			// Render currentLine
			const lineMetrics = baseCtx.measureText(currentLine)
			const textWidthPx = lineMetrics.width

			const textCanvas = document.createElement('canvas')
			textCanvas.width = textWidthPx * scaleFactor
			textCanvas.height = fontSize * 2 * scaleFactor
			const textCtx = textCanvas.getContext('2d')!
			textCtx.scale(scaleFactor, scaleFactor)
			textCtx.font = `${fontStyle} ${fontSize}pt Helvetica`
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

	// Render leftover text
	if (currentLine.trim()) {
		const lineMetrics = baseCtx.measureText(currentLine)
		const textWidthPx = lineMetrics.width

		const textCanvas = document.createElement('canvas')
		textCanvas.width = textWidthPx * scaleFactor
		textCanvas.height = fontSize * 2 * scaleFactor
		const textCtx = textCanvas.getContext('2d')!
		textCtx.scale(scaleFactor, scaleFactor)
		textCtx.font = `${fontStyle} ${fontSize}pt Helvetica`
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
	 * Generate the PDF version.
	 * This is for normal prints (when SAVE_ON_PRINT = false).
	 */
	async function generatePdfFromSuborder(
		order: OrderDTO,
		suborder: SuborderDTO,
		index: number,
		options: QRPrintOptions = {},
	): Promise<Blob> {
		const { dpi = 203, labelWidthMm = 80, labelHeightMm = 80 } = options
		const pxPerMm = dpi / 25.4

		const qrValue = generateSubOrderQR(suborder)

		const line1 = `#${order.displayNumber} ${order.customerName} (${index + 1}/${order.subOrders.length})`
		const line2 = `${suborder.productSize.productName} ${suborder.productSize.sizeName} (${suborder.productSize.size} ${suborder.productSize.unit.name.toLowerCase()}, ${MACHINE_CATEGORY_FORMATTED[suborder.productSize.machineCategory].toLowerCase()})`

		// Draw QR code on a canvas
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
		const qrScaleFactor = 0.9
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

		const qrDataURL = canvas.toDataURL('image/png', 1.0)

		// Create jsPDF doc
		const doc = new jsPDF({
			orientation: labelWidthMm > labelHeightMm ? 'landscape' : 'portrait',
			unit: 'mm',
			format: [labelWidthMm, labelHeightMm],
		})

		const baseFontSize = Math.min(10, labelWidthMm / 8)
		const margin = Math.max(2, labelWidthMm * 0.025)
		let currentY = margin
		const maxTextWidth = labelWidthMm - margin * 2

		// Draw text lines
		currentY = drawText(doc, line1, margin, currentY, maxTextWidth, baseFontSize, true)
		currentY = drawText(doc, line2, margin, currentY, maxTextWidth, baseFontSize * 0.9)

		// Then place the QR code
		const qrSpacing = baseFontSize * 0.5
		currentY += qrSpacing
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

		return doc.output('blob')
	}

	/**
	 * Convert text to CP1251 byte representation for ZPL
	 */
	function textToCP1251Bytes(input: string): number[] {
		return Array.from(input).map(char => {
			const code = char.charCodeAt(0)
			// If the character is in the Cyrillic block, adjust it.
			if (code >= 0x0410 && code <= 0x044f) {
				// CP1251 for Cyrillic (simplified mapping)
				return code - 0x350
			}
			// For basic ASCII, the byte is the same:
			return code < 128 ? code : 63 // unknown chars become '?'
		})
	}

	/**
	 * Escape text for ZPL format with CP1251 encoding
	 */
	function escapeZplCP1251(input: string): string {
		const bytes = textToCP1251Bytes(input)
		return bytes.map(byte => `\\${byte.toString(16).padStart(2, '0')}`).join('')
	}

	/**
	 * Improved and simplified PRN generator with better layout
	 * for different label sizes, especially 60x40mm
	 */
	async function generatePrnFromSuborder(
		order: OrderDTO,
		suborder: SuborderDTO,
		index: number,
		options: QRPrintOptions = {},
	) {
		// Default label dimensions in mm
		const { labelWidthMm = 80, labelHeightMm = 80 } = options

		// Convert mm to dots (8 dots per mm for 203 DPI printers)
		const labelWidthDots = Math.round(labelWidthMm * 8)
		const labelHeightDots = Math.round(labelHeightMm * 8)

		// Calculate fixed margins - use a percentage of the label dimensions
		// For smaller labels, we use a smaller percentage to maximize usable space
		const marginX = Math.round(labelWidthDots * 0.02) // 2% margin
		const textAreaWidth = labelWidthDots - marginX * 2

		// Generate the order content
		const qrValue = generateSubOrderQR(suborder)
		const line1 = `#${order.displayNumber} ${order.customerName} (${index + 1}/${order.subOrders.length})`
		const line2 = `${suborder.productSize.productName} ${suborder.productSize.sizeName} (${suborder.productSize.size} ${suborder.productSize.unit.name.toLowerCase()}, ${MACHINE_CATEGORY_FORMATTED[suborder.productSize.machineCategory].toLowerCase()})`

		// Escape text for ZPL encoding
		const escapedLine1 = escapeZplCP1251(line1)
		const escapedLine2 = escapeZplCP1251(line2)

		// Calculate font sizes based on label width
		// For different size labels, fonts should scale proportionally
		// with minimum sizes to ensure readability
		const fontScale = Math.min(1, Math.max(0.6, labelWidthMm / 80))

		// For 80mm width, we use 38pt for line1 and 29pt for line2
		// Scale these values for different widths while keeping minimums
		const line1FontSize = Math.max(20, Math.round(38 * fontScale))
		const line2FontSize = Math.max(16, Math.round(29 * fontScale))

		// Determine vertical spacing
		// The original static values seem to work well, so adapt them for different heights
		// For 80mm height: line1Y = 16, line2Y = 55
		const heightScale = Math.min(1, Math.max(0.5, labelHeightMm / 80))
		const line1Y = Math.max(10, Math.round(16 * heightScale))
		const line2Y = Math.max(line1Y + line1FontSize + 5, Math.round(55 * heightScale))

		// Calculate QR code position and size
		// For 80mm labels, QR starts at ~144,105
		// Target about 33% of height down for QR code Y position
		const qrY = Math.round(labelHeightDots * 0.33)

		// Center the QR code horizontally
		// For wide layouts (width > height), make QR code about 50% of height
		// For square/tall layouts, make QR code about a comfortable size
		const qrSize = Math.round(Math.min(labelWidthDots * 0.4, labelHeightDots * 0.5))
		const qrX = Math.round((labelWidthDots - qrSize) / 2)

		// QR magnification factor - higher gives larger modules
		// For 80mm labels, 10 works well. Scale for different sizes.
		const qrMagMin = labelHeightMm <= 40 ? 6 : 8 // Ensure readability on small labels
		const qrMagnification = Math.max(qrMagMin, Math.round(qrSize / 30))

		// Build the ZPL commands
		const zplCommands = `^XA
  ^CI33
  ^PW${labelWidthDots}
  ^LL${labelHeightDots}

  ^FO${marginX},${line1Y}
  ^FB${textAreaWidth},2,0,C,0
  ^A0N,${line1FontSize},${Math.round(line1FontSize * 0.8)}^FH\\
  ^FD${escapedLine1}^FS

  ^FO${marginX},${line2Y}
  ^FB${textAreaWidth},3,0,C,0
  ^A0N,${line2FontSize},${Math.round(line2FontSize * 0.8)}^FH\\
  ^FD${escapedLine2}^FS

  ^FO${qrX},${qrY}
  ^BQN,2,${qrMagnification}
  ^FDLA,${qrValue}^FS

  ^MMC
  ^XZ`.trim()

		return new Blob([zplCommands], { type: 'application/octet-stream' })
	}

	/**
	 * Decide which type of label to generate (PDF vs. raw spool .prn).
	 */
	async function generateQRFromSuborder(
		order: OrderDTO,
		suborder: SuborderDTO,
		index: number,
		options: QRPrintOptions = {},
	): Promise<Blob> {
		if (SAVE_ON_PRINT) {
			// Generate raw spool file (.prn) with ZPL commands
			return generatePrnFromSuborder(order, suborder, index, options)
		} else {
			// Generate PDF
			return generatePdfFromSuborder(order, suborder, index, options)
		}
	}

	return {
		generateQRFromSuborder,
	}
}

export function useOrderQRPrinter() {
	const { print } = usePrinter()
	const { generateQRFromSuborder } = useOrderGenerateQR()

	/**
	 * Prints (or saves) one or more QR labels for a suborder.
	 * If SAVE_ON_PRINT = true, we generate a .prn spool file
	 * (ZPL commands, cyrillic in hex).
	 * Otherwise, we generate a PDF and pass it to the default print logic.
	 */
	async function printSubOrderQR(
		order: OrderDTO,
		subOrders: SuborderDTO | SuborderDTO[],
		options?: QRPrintOptions,
	) {
		try {
			const suborders = Array.isArray(subOrders) ? subOrders : [subOrders]
			// Generate a Blob for each suborder (PDF or .prn)
			const labelBlobs = await Promise.all(
				suborders.map((sb, idx) => {
					const printIndex =
						suborders.length === 1 ? order.subOrders.findIndex(s => s.id === sb.id) : idx
					return generateQRFromSuborder(order, sb, printIndex, options)
				}),
			)

			// Hand them off to the printer or saving logic
			await print(labelBlobs, options)
		} catch (error) {
			console.error('Error generating or printing QR code:', error)
			throw error
		}
	}

	return {
		printSubOrderQR,
	}
}
