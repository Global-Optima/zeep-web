// tests/phone.utils.spec.ts
import { formatPhoneNumber } from './fomat-phone-number.utils'
import { describe, expect, it } from 'vitest'

describe('formatPhoneNumber', () => {
	it('should format a valid +7‑number correctly', () => {
		expect(formatPhoneNumber('+79667778899')).toBe('+7 966 777-88-99')
	})

	it('should return empty string for missing “+7” prefix', () => {
		expect(formatPhoneNumber('96677778899')).toBe('')
	})

	it('should return empty string for too-short numbers', () => {
		expect(formatPhoneNumber('+7123456')).toBe('')
	})

	it('should return empty string for too-long numbers', () => {
		expect(formatPhoneNumber('+796677788990')).toBe('')
	})

	it('should return empty string for completely invalid strings', () => {
		expect(formatPhoneNumber('invalid')).toBe('')
		expect(formatPhoneNumber('')).toBe('')
	})
})
