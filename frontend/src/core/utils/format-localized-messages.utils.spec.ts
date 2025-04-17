import { describe, expect, it } from 'vitest'
import { formatLocalizedMessage, sanitizeString } from './format-localized-messages.utils'

describe('sanitizeString', () => {
	it('returns the same text for safe strings', () => {
		expect(sanitizeString('Hello, world!')).toBe('Hello, world!')
	})

	it('allows default safe HTML tags', () => {
		// <b> is in xss default whitelist
		expect(sanitizeString('<b>bold</b>')).toBe('<b>bold</b>')
	})

	it('strips disallowed <script> tags and content', () => {
		expect(sanitizeString('<script>alert(1)</script>hello')).toBe('hello')
	})

	it('strips all tags if whiteList is empty', () => {
		const raw = '<div>content</div>'
		expect(sanitizeString(raw, { whiteList: {} })).toBe('content')
	})

	it('escapes malicious attributes', () => {
		const raw = '<img src="x" onerror="alert(1)" />'
		expect(sanitizeString(raw)).toBe('&lt;img src="x" onerror="alert(1)" /&gt;')
	})
})

describe('formatLocalizedMessage', () => {
	it('wraps *text* in a span and leaves other text intact', () => {
		const input = 'Hello *world*!'
		expect(formatLocalizedMessage(input)).toBe('Hello <span class="font-medium">"world"</span>!')
	})

	it('replaces newlines with <br> tags', () => {
		const input = 'Line1\nLine2'
		expect(formatLocalizedMessage(input)).toBe('Line1<br>Line2')
	})

	it('sanitizes inner content of *...* before wrapping', () => {
		const input = '*<b>bad</b>*'
		// inner <b> should be stripped by sanitize
		expect(formatLocalizedMessage(input)).toBe('<span class="font-medium">"bad"</span>')
	})

	it('escapes disallowed tags after formatting', () => {
		const input = '*safe*<script>alert</script>'
		const output = formatLocalizedMessage(input)
		expect(output).toContain('<span class="font-medium">"safe"</span>')
		// script tags should be escaped
		expect(output).toContain('&lt;script&gt;alert&lt;/script&gt;')
	})

	it('combines markdown and newlines and sanitizes all', () => {
		const input = '*one*\n<script>evil()</script>'
		const output = formatLocalizedMessage(input)
		expect(output).toBe(
			'<span class="font-medium">"one"</span><br>&lt;script&gt;evil()&lt;/script&gt;',
		)
	})
})
