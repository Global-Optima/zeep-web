// src/core/utils/string.utils.ts
import xss, { whiteList as xssWhiteList, type IFilterXSSOptions } from 'xss'

/**
 * Sanitizes a given string using xss, with these rules:
 * 1. Remove <script> tags and their content entirely.
 * 2. If caller passes a whiteList, strip all other tags (keep inner text).
 * 3. Otherwise, allow default xssWhiteList (minus <img>), and escape every other tag.
 *
 * @param str - The input string to sanitize.
 * @param options - Optional xss options to customize sanitization.
 * @returns The sanitized string.
 */
export const sanitizeString = (str: string, options?: IFilterXSSOptions): string => {
	// 1) Strip out any <script>…</script> blocks
	const noScript = str.replace(/<script[\s\S]*?<\/script>/gi, '')

	// 2) If caller provided a whiteList, strip (remove) all tags not in it
	if (options?.whiteList !== undefined) {
		return xss(noScript, {
			whiteList: options.whiteList,
			stripIgnoreTag: true,
		})
	}

	// 3) Default behavior: allow xss's default tags (minus <img>), escape everything else
	const defaultWhiteList = { ...xssWhiteList }
	delete defaultWhiteList.img

	return xss(noScript, {
		whiteList: defaultWhiteList,
		stripIgnoreTag: false, // escape disallowed tags rather than stripping
	})
}

/**
 * Formats a localized message:
 * 1. Transforms *markdown* into <span class="font-medium">"…"</span>.
 * 2. Converts newlines to <br>.
 * 3. Sanitizes final HTML to allow only <span class="…"> and <br>, escaping everything else.
 *
 * @param s - The raw localized string.
 * @returns A safe HTML string for v-html.
 */
export const formatLocalizedMessage = (s: string): string => {
	// 1) Markdown‐style emphasis
	const withSpans = s.replace(/\*(.*?)\*/g, (_m, content) => {
		// Strip any tags inside the markdown span
		const inner = sanitizeString(content, { whiteList: {} })
		return `<span class="font-medium">"${inner}"</span>`
	})

	// 2) Newlines → <br>
	const withBreaks = withSpans.replace(/\n/g, '<br>')

	// 3) Final sanitize: only allow <span class="…"> and <br>, all else escaped
	return xss(withBreaks, {
		whiteList: { span: ['class'], br: [] },
		stripIgnoreTag: false, // ensure disallowed tags (like <script>) get escaped
	})
}
