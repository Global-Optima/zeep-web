import { SAVED_BARISTA_QR_SETTINGS_KEY } from '@/core/constants/qr.constants'
import type { SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { jsPDF } from 'jspdf'
import QRCode from 'qrcode-generator'
import { usePrinter, type PrintOptions } from './use-print.hook'
import { MachineCategory } from '@/modules/admin/product-categories/utils/category-options'

const categoryEmojis: Record<MachineCategory, string> = {
  [MachineCategory.TEA]: 'Чай',
  [MachineCategory.COFFEE]: 'Кофе',
  [MachineCategory.ICE_CREAM]: 'Мороженое'
};

export interface QRPrintOptions extends PrintOptions {
  labelWidthMm?: number
  labelHeightMm?: number
  dpi?: number
  textContent?: {
    customerNameField?: string
    productSizeField?: string
    additivesField?: string
  }
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
   * Generates a QR code with optional text content and embeds it into a PDF.
   * @param qrValue - The string value to encode in the QR code.
   * @param options - Configuration for QR printing including optional text content.
   */
  async function generateQR(
    qrValue: string,
    options: QRPrintOptions = {}
  ): Promise<Blob> {
    const { dpi = 203, labelWidthMm = 100, labelHeightMm = 100, textContent } = options

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
    // 3) Determine layout - text area and QR area
    // -------------------------------
    const hasTextContent = textContent && (textContent.customerNameField || textContent.productSizeField || textContent.additivesField)
    const textAreaHeightPx = hasTextContent ? Math.round(usableHeightPx * 0.3) : 0
    const qrAreaHeightPx = usableHeightPx - textAreaHeightPx

    // Calculate QR code size - make it slightly smaller than the QR area
    const qrSizePx = Math.min(usableWidthPx * 0.9, qrAreaHeightPx * 0.9)

    // Center QR code horizontally and position it in the QR area
    const qrMarginHorizontalPx = (usableWidthPx - qrSizePx) / 2
    const qrMarginVerticalPx = textAreaHeightPx + (qrAreaHeightPx - qrSizePx) / 2

    // -------------------------------
    // 4) Draw text information if provided
    // -------------------------------
    if (hasTextContent) {
      ctx.fillStyle = '#000000'

      // Calculate maximum available width
      const maxTextWidth = usableWidthPx * 0.9

      // Base font sizes
      const baseFontSizes = {
        customerName: Math.max(20, Math.min(30, labelWidthMm / 3)),
        productSize: Math.max(20, Math.min(30, labelWidthMm / 3)),
        additives: Math.max(15, Math.min(25, labelWidthMm / 3)) * 0.8
      }

      // Calculate total text height needed
      const textBlocks: Array<{
        text: string;
        initialSize: number;
        weight: string;
        lineSpacing?: number;
      }> = [];

      // Prepare text blocks and calculate total height
      if (textContent.customerNameField) {
        textBlocks.push({
          text: textContent.customerNameField,
          initialSize: baseFontSizes.customerName,
          weight: '500',
          lineSpacing: 1.2
        });
      }

      if (textContent.productSizeField) {
        textBlocks.push({
          text: textContent.productSizeField,
          initialSize: baseFontSizes.productSize,
          weight: '400',
          lineSpacing: 1.2
        });
      }

      if (textContent.additivesField) {
        textBlocks.push({
          text: textContent.additivesField,
          initialSize: baseFontSizes.additives,
          weight: '400',
          lineSpacing: 0.8
        });
      }

      // First pass: calculate total height needed
      let calculatedHeight = 0;
      const textMeasurements = textBlocks.map(block => {
        const { fontSize, lines } = getOptimalTextFit(
          ctx,
          block.text,
          maxTextWidth,
          block.initialSize,
          block.weight
        );

        const blockHeight = lines.length * fontSize * (block.lineSpacing || 1.2);
        calculatedHeight += blockHeight;

        return {
          ...block,
          fontSize,
          lines,
          blockHeight
        };
      });

      // Calculate available text area (30% of total height)
      const textAreaHeightPx = usableHeightPx * 0.3;

      // Calculate vertical centering
      const verticalMargin = Math.max(0, (textAreaHeightPx - calculatedHeight) / 2);
      let currentY = verticalMargin;

      // Second pass: actual drawing
      for (const block of textMeasurements) {
        ctx.font = `${block.weight} ${block.fontSize}px Arial, sans-serif`;
        ctx.textAlign = 'center';

        for (const line of block.lines) {
          ctx.fillText(line, usableWidthPx / 2, currentY + block.fontSize);
          currentY += block.fontSize * (block.lineSpacing || 1.2);
        }
      }
    }

    // -------------------------------
    // 5) Generate and draw QR Code
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
    doc.addImage(dataURL, 'PNG', 0, 0, labelWidthMm, labelHeightMm)

    // Return as a Blob for printing
    return doc.output('blob')
  }

  // Helper function to calculate optimal font size and line breaks
  function getOptimalTextFit(
    ctx: CanvasRenderingContext2D,
    text: string,
    maxWidth: number,
    initialSize: number,
    weight: string = '400',
    maxLines: number = 3
  ): { fontSize: number; lines: string[] } {
    let fontSize = initialSize
    let lines: string[] = []

    // First try with initial font size
    ctx.font = `${weight} ${fontSize}px Arial, sans-serif`

    // Check if text fits on one line
    const metrics = ctx.measureText(text)
    if (metrics.width <= maxWidth) {
      return { fontSize, lines: [text] }
    }

    // If not, try reducing font size
    while (fontSize > 10) {
      ctx.font = `${weight} ${fontSize}px Arial, sans-serif`
      lines = wrapTextToLines(ctx, text, maxWidth, maxLines)

      if (lines.length <= maxLines) {
        break
      }

      fontSize -= 1 // Reduce font size and try again
    }

    return { fontSize, lines }
  }

  // Helper function to wrap text into multiple lines (handle long unbroken strings)
  function wrapTextToLines(
    ctx: CanvasRenderingContext2D,
    text: string,
    maxWidth: number,
    maxLines: number = 3
  ): string[] {
    const lines: string[] = [];
    let currentLine = '';
    let currentLineWidth = 0;

    // Function to add a line and reset counters
    const addLine = (line: string) => {
      lines.push(line);
      currentLine = '';
      currentLineWidth = 0;
    };

    // Split into words first
    const words = text.split(' ');

    for (const word of words) {
      const wordWidth = ctx.measureText(word).width;

      // If word fits in current line
      if (currentLineWidth + wordWidth <= maxWidth) {
        currentLine += (currentLine ? ' ' : '') + word;
        currentLineWidth += wordWidth + (currentLine ? ctx.measureText(' ').width : 0);
      }
      // If word doesn't fit but is too long for a line by itself
      else if (wordWidth > maxWidth) {
        // If we already have content in current line, push it first
        if (currentLine) {
          addLine(currentLine);
          if (lines.length >= maxLines) break;
        }

        // Handle the long word by breaking it into parts
        let remainingWord = word;
        while (remainingWord) {
          // Find how many characters fit
          let charsThatFit = 0;
          let testWidth = 0;

          while (charsThatFit < remainingWord.length) {
            const testString = remainingWord.substring(0, charsThatFit + 1);
            testWidth = ctx.measureText(testString).width;
            if (testWidth > maxWidth) break;
            charsThatFit++;
          }

          // Add the part that fits
          const part = remainingWord.substring(0, charsThatFit);
          addLine(part);
          if (lines.length >= maxLines) break;

          // Prepare remaining part
          remainingWord = remainingWord.substring(charsThatFit);
        }

        if (lines.length >= maxLines) break;
      }
      else {
        // Word doesn't fit in current line but would fit in new line
        if (currentLine) {
          addLine(currentLine);
          if (lines.length >= maxLines) break;
        }
        currentLine = word;
        currentLineWidth = wordWidth;
      }
    }

    // Add any remaining content
    if (currentLine && lines.length < maxLines) {
      addLine(currentLine);
    }

    // Add ellipsis to last line if we truncated
    if (lines.length === maxLines) {
      const lastLine = lines[lines.length - 1];
      const ellipsis = '...';
      const ellipsisWidth = ctx.measureText(ellipsis).width;

      // Remove characters until ellipsis fits
      let trimmedLine = lastLine;
      while (trimmedLine && ctx.measureText(trimmedLine + ellipsis).width > maxWidth) {
        trimmedLine = trimmedLine.substring(0, trimmedLine.length - 1);
      }

      lines[lines.length - 1] = trimmedLine + (trimmedLine.length < lastLine.length ? ellipsis : '');
    }

    return lines;
  }

  return { generateQR }
}

export function useQRPrinter() {
  const { print } = usePrinter()
  const { generateQR } = useGenerateQR()

  /**
   * Generic QR printing function that can be used for any QR code printing
   * @param qrValue - The encoded QR value or array of values
   * @param options - Optional print configurations
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

  /**
   * Order-specific QR printing function (maintained for backward compatibility)
   * @param subOrders - Array of suborders to print QR codes for
   * @param orderNumber - Order number to display on each QR
   * @param customerName - Customer name to display on each QR
   * @param options - Optional print configurations
   */
  const printOrderQR = async (
    subOrders: SuborderDTO | SuborderDTO[],
    orderNumber?: number,
    customerName?: string,
    options?: QRPrintOptions
  ) => {
    const ordersArray = Array.isArray(subOrders) ? subOrders : [subOrders]

    const qrValues = ordersArray.map(subOrder => {
      const qrValue = generateSubOrderQR(subOrder)

      const machineCategoryText = categoryEmojis[subOrder.productSize?.machineCategory as MachineCategory] ?? '';

      const productSizeField = subOrder.productSize?.productName && subOrder.productSize?.sizeName
        ? `${subOrder.productSize?.productName} ${subOrder.productSize?.sizeName} (${machineCategoryText}, ${subOrder.productSize?.size} ${subOrder.productSize?.unit?.name})`
        : undefined

      const additivesField = subOrder.additives?.length
        ? `${subOrder.additives.map(a => a.additive.name).join(', ')}`
        : undefined

      return {
        qrValue,
        options: {
          ...options,
          textContent: {
            customerNameField: orderNumber + " " + customerName,
            productSizeField,
            additivesField
          }
        }
      }
    })

    // Use the generic printQR function internally
    await printQR(
      qrValues.map(v => v.qrValue),
      {
        ...options,
        textContent: qrValues[0]?.options.textContent
      }
    )
  }

  return {
    printQR,       // Generic version
    printOrderQR   // Order-specific version
  }
}
