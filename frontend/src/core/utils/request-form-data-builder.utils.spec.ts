import { describe, expect, it } from 'vitest'
import { buildFormData } from './request-form-data-builder.utils'

describe('buildFormData', () => {
	it('appends strings, numbers, and booleans correctly', () => {
		const dto = { a: 'hello', b: 42, c: true, d: null, e: undefined }
		const fd = buildFormData(dto)
		const entries = Object.fromEntries(fd.entries())

		expect(entries).toHaveProperty('a', 'hello')
		expect(entries).toHaveProperty('b', '42')
		expect(entries).toHaveProperty('c', 'true')
		expect(entries).not.toHaveProperty('d')
		expect(entries).not.toHaveProperty('e')
	})

	it('serializes arrays as JSON strings', () => {
		const arr = [1, 2, 3]
		const dto = { list: arr }
		const fd = buildFormData(dto)
		expect(fd.get('list')).toBe(JSON.stringify(arr))
	})

	it('recursively flattens nested objects', () => {
		const dto = { obj: { x: 'X', y: 99 } }
		const fd = buildFormData(dto)
		expect(fd.get('obj[x]')).toBe('X')
		expect(fd.get('obj[y]')).toBe('99')
	})

	it('flattens deeply nested objects', () => {
		const dto = { a: { b: { c: 'deep' } } }
		const fd = buildFormData(dto)
		expect(fd.get('a[b][c]')).toBe('deep')
	})

	it('appends File instances directly', () => {
		const file = new File(['content'], 'test.txt', { type: 'text/plain' })
		const dto = { file }
		const fd = buildFormData(dto)
		const entry = fd.get('file')
		expect(entry).toBeInstanceOf(File)
		if (entry instanceof File) {
			expect(entry.name).toBe('test.txt')
			expect(entry.type).toBe('text/plain')
		}
	})

	it('serializes array-of-objects as JSON strings', () => {
		const arrObj = [{ a: 1 }, { b: 2 }]
		const dto = { arrObj }
		const fd = buildFormData(dto)
		expect(fd.get('arrObj')).toBe(JSON.stringify(arrObj))
	})
})
