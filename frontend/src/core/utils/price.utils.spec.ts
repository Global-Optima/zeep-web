import { describe, expect, it } from 'vitest'
import { formatPrice } from './price.utils'

describe('formatPrice', () => {
	it('should format a valid number correctly', () => {
		const result = formatPrice(123456.789)
		expect(result).toBe('123 456,789 ₸')
	})

	it('should format a valid numeric string correctly', () => {
		const result = formatPrice('123456.789')
		expect(result).toBe('123 456,789 ₸')
	})

	it('should throw an error for an invalid numeric string', () => {
		expect(() => formatPrice('invalid')).toThrowError(
			'Invalid value provided. Value must be a valid number or numeric string.',
		)
	})

	it('should throw an error for a non-numeric value', () => {
		expect(() => formatPrice(NaN)).toThrowError(
			'Invalid value provided. Value must be a valid number or numeric string.',
		)
	})

	it('should handle zero correctly', () => {
		const result = formatPrice(0)
		expect(result).toBe('0 ₸')
	})

	it('should handle negative numbers correctly', () => {
		const result = formatPrice(-1234.56)
		expect(result).toBe('-1 234,56 ₸')
	})

	it('should handle large numbers correctly', () => {
		const result = formatPrice(123456789012.34)
		expect(result).toBe('123 456 789 012,34 ₸')
	})
})
