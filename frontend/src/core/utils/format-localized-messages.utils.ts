import xss, { type IFilterXSSOptions } from 'xss'

/**
 * Sanitizes a given string using xss with provided options.
 *
 * @param s - The input string to sanitize.
 * @param options - Optional xss options to customize sanitization.
 * @returns The sanitized string.
 */
export const sanitizeString = (s: string, options?: IFilterXSSOptions): string => {
	return xss(s, options)
}

/**
 * Formats a localized message with custom styling.
 *
 * This function transforms markdown-style asterisks into a <span> element,
 * replaces newlines with <br>, and sanitizes both the inner content and the final HTML.
 *
 * @param s - The raw localized string.
 * @returns A sanitized HTML string safe for v-html binding.
 */
export const formatLocalizedMessage = (s: string): string => {
	// Replace markdown-style *text* with a <span> wrapping the sanitized text.
	const formatted = s
		.replace(/\*(.*?)\*/g, (_match, content) => {
			// Sanitize the inner content, allowing no HTML tags.
			const safeContent = sanitizeString(content, { whiteList: {} })
			return `<span class="font-medium">"${safeContent}"</span>`
		})
		// Replace newline characters with <br> tags.
		.replace(/\n/g, '<br>')

	// Sanitize the entire formatted string, allowing only <span> tags (with class attribute) and <br> tags.
	return sanitizeString(formatted, {
		whiteList: {
			span: ['class'],
			br: [],
		},
	})
}
