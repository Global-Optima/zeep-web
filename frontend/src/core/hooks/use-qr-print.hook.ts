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
   * Generates a QR code with human-readable text and embeds it into a PDF (correctly scaled for printing).
   * @param qrValue - The string value to encode in the QR code.
   * @param subOrder - The SuborderDTO to extract display information from.
   * @param customerName - The name of the customer who placed the order.
   * @param machineCategory - The machine category of the suborder product
   * @param options - Configuration for QR printing.
   */
  async function generateQR(
    qrValue: string,
    subOrder?: SuborderDTO,
    customerName?: string,
    options: QRPrintOptions = {}
  ): Promise<Blob> {
    const { dpi = 203, labelWidthMm = 100, labelHeightMm = 100 } = options

    // -------------------------------
    // 1) Define Label Size in mm
    // -------------------------------
    const pxPerMm = dpi / 25.4 // Conversion factor: DPI to pixels per millimeter

    // Compute usable area inside the label
    const usableWidthPx = Math.round(labelWidthMm * pxPerMm)
    const usableHeightPx = Math.round(labelHeightMm * pxPerMm)

    // Create a high-res canvas
    const canvas = document.createElement('canvas')
    canvas.width = usableWidthPx
    canvas.height = usableHeightPx

    // -------------------------------
    // 2) Clear canvas with white background
    // -------------------------------
    const ctx = canvas.getContext('2d')!
    ctx.fillStyle = '#FFFFFF'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // -------------------------------
    // 3) Determine layout - explicit text area and QR area
    // -------------------------------
    // Allocate top 30% for text, bottom 70% for QR
    const textAreaHeightPx = Math.round(usableHeightPx * 0.3)
    const qrAreaHeightPx = usableHeightPx - textAreaHeightPx

    // Calculate QR code size - make it slightly smaller than the QR area
    const qrSizePx = Math.min(usableWidthPx * 0.9, qrAreaHeightPx * 0.9)

    // Center QR code horizontally and position it in the QR area
    const qrMarginHorizontalPx = (usableWidthPx - qrSizePx) / 2
    const qrMarginVerticalPx = textAreaHeightPx + (qrAreaHeightPx - qrSizePx) / 2

    // -------------------------------
    // 4) Draw text information in the text area
    // -------------------------------
    if (subOrder) {
      ctx.fillStyle = '#000000'

      // Base font size on label width (responsive)
      const baseFontSize = Math.max(50, Math.min(35, labelWidthMm / 3))
      const lineSpacing = baseFontSize * 1.2

      // Start position for text (top of text area with margin)
      let currentY = lineSpacing

      // 1. Draw customer name (larger font)
      ctx.font = `bold ${baseFontSize + 4}px Arial, sans-serif`
      ctx.textAlign = 'center'

      if (customerName) {
        ctx.fillText(customerName, usableWidthPx / 2, currentY)
        currentY += lineSpacing
      }

      // 2. Draw product name with category emoji and size
      if (subOrder.productSize?.productName && subOrder.productSize?.sizeName) {
        ctx.font = `bold ${baseFontSize}px Arial, sans-serif`;

        // Combine product name, and size
        const productText = `${subOrder.productSize.productName} - ${subOrder.productSize.sizeName}`;

        // Centered text
        ctx.fillText(productText, usableWidthPx / 2, currentY);

        currentY += lineSpacing; // Move to the next line
      }

      // 4. Draw additives if any
      if (subOrder.additives?.length) {
        ctx.font = `${baseFontSize - 2}px Arial, sans-serif`

        // Draw "With:" text
        ctx.fillText("Добавки:", usableWidthPx / 2, currentY);
        currentY += lineSpacing * 0.8;

        // Join additive names with commas for compact display
        const additiveNames = subOrder.additives
          .map(a => a.additive.name)
          .join(", ");

        // Handle text wrapping for additives
        const maxWidth = usableWidthPx * 0.9;
        wrapText(ctx, additiveNames, usableWidthPx / 2, currentY, maxWidth, lineSpacing * 0.8);
      }

      // Draw separator line between text and QR
      ctx.strokeStyle = '#CCCCCC';
      ctx.lineWidth = 1;
      ctx.beginPath();
      ctx.moveTo(10, textAreaHeightPx - 5);
      ctx.lineTo(usableWidthPx - 10, textAreaHeightPx - 5);
      ctx.stroke();
    }

    // -------------------------------
    // 5) Generate and draw QR Code in the QR area
    // -------------------------------
    const qr = QRCode(0, 'M') // M = Medium Error Correction
    qr.addData(qrValue)
    qr.make()

    ctx.fillStyle = '#000000' // Black for QR Code
    const qrModules = qr.getModuleCount()
    const cellSize = qrSizePx / qrModules

    for (let row = 0; row < qrModules; row++) {
      for (let col = 0; col < qrModules; col++) {
        if (qr.isDark(row, col)) {
          ctx.fillRect(
            qrMarginHorizontalPx + col * cellSize,
            qrMarginVerticalPx + row * cellSize,
            cellSize,
            cellSize
          )
        }
      }
    }

    // -------------------------------
    // 6) Create PDF with the combined image
    // -------------------------------
    const isLandscape = labelWidthMm > labelHeightMm
    const doc = new jsPDF({
      orientation: isLandscape ? 'landscape' : 'portrait',
      unit: 'mm',
      format: [labelWidthMm, labelHeightMm],
    })

    // Convert the canvas to data URL
    const dataURL = canvas.toDataURL('image/png', 100)

    // Insert the image into the PDF (full size)
    const xStartMm = 0
    const yStartMm = 0
    doc.addImage(dataURL, 'PNG', xStartMm, yStartMm, labelWidthMm, labelHeightMm)

    // Return as a Blob for printing
    return doc.output('blob')
  }

  // Helper function to wrap text
  function wrapText(ctx, text, x, y, maxWidth, lineHeight) {
    const words = text.split(' ');
    let line = '';
    let testLine = '';
    let lineCount = 0;

    for (let n = 0; n < words.length; n++) {
      testLine = line + words[n] + ' ';
      const metrics = ctx.measureText(testLine);
      const testWidth = metrics.width;

      if (testWidth > maxWidth && n > 0) {
        ctx.fillText(line, x, y);
        line = words[n] + ' ';
        y += lineHeight;
        lineCount++;

        // Limit to max 2 lines for additives
        if (lineCount >= 2) {
          // Add ellipsis if we're cutting off text
          if (n < words.length - 1) {
            // Remove last word and add ellipsis
            line = line.substring(0, line.lastIndexOf(' ')) + '...';
          }
          ctx.fillText(line, x, y);
          break;
        }
      } else {
        line = testLine;
      }
    }

    // Draw the last line if we haven't hit the limit
    if (lineCount < 2) {
      ctx.fillText(line, x, y);
    }
  }


  return { generateQR }
}

export function useQRPrinter() {
  const { print } = usePrinter()
  const { generateQR } = useGenerateQR()

  /**
   * Print QR codes for suborders with human-readable information
   * @param subOrders - Array of suborders to print QR codes for
   * @param customerName - Customer name to display on each QR
   * @param machineCategory - The machine category of the suborder product
   * @param options - Optional print configurations
   */
  const printQR = async (
    subOrders: SuborderDTO[],
    customerName?: string,
    options?: QRPrintOptions
  ) => {
    try {
      const qrPdfBlobs = await Promise.all(
        subOrders.map(subOrder => {
          const qrValue = generateSubOrderQR(subOrder)
          return generateQR(qrValue, subOrder, customerName, options)
        })
      )

      await print(qrPdfBlobs, options)
    } catch (error) {
      console.error('Error generating or printing QR code:', error)
      throw error
    }
  }

  // For backward compatibility with original API
  const printQRValues = async (qrValues: string | string[], options?: QRPrintOptions) => {
    try {
      const values = Array.isArray(qrValues) ? qrValues : [qrValues]
      const qrPdfBlobs = await Promise.all(values.map(v => generateQR(v, undefined, undefined, options)))
      await print(qrPdfBlobs, options)
    } catch (error) {
      console.error('Error generating or printing QR code:', error)
      throw error
    }
  }

  return { printQR, printQRValues }
}
