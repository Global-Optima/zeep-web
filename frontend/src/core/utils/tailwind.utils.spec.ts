import { describe, expect, it } from 'vitest'
import { cn } from './tailwind.utils'

describe('cn', () => {
	it('should merge single class name correctly', () => {
		const result = cn('bg-red-500')
		expect(result).toBe('bg-red-500')
	})

	it('should merge multiple class names correctly', () => {
		const result = cn('bg-red-500', 'text-white')
		expect(result).toBe('bg-red-500 text-white')
	})

	it('should handle conditional class names correctly', () => {
		const isActive = true
		const result = cn('bg-red-500', isActive && 'text-white', 'p-4')
		expect(result).toBe('bg-red-500 text-white p-4')
	})

	it('should ignore falsey values', () => {
		const result = cn('bg-red-500', false && 'text-white', null, undefined, 'p-4')
		expect(result).toBe('bg-red-500 p-4')
	})

	it('should handle an array of class names', () => {
		const result = cn(['bg-red-500', 'text-white'])
		expect(result).toBe('bg-red-500 text-white')
	})

	it('should merge Tailwind classes correctly with overrides', () => {
		const result = cn('bg-red-500', 'bg-blue-500')
		expect(result).toBe('bg-blue-500')
	})

	it('should handle an empty input gracefully', () => {
		const result = cn()
		expect(result).toBe('')
	})
})
