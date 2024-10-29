import { describe, expect, it } from 'vitest'

function add(a: number, b: number) {
	return a + b
}

function subtract(a: number, b: number) {
	return a - b
}

function multiply(a: number, b: number) {
	return a * b
}

function divide(a: number, b: number) {
	if (b === 0) {
		return { quotient: 0, remainder: 0 }
	}
	return { quotient: Math.floor(a / b), remainder: a % b }
}

describe('Math Functions', () => {
	it('should add two numbers correctly', () => {
		expect(add(2, 3)).toBe(5)
	})

	it('should subtract two numbers correctly', () => {
		expect(subtract(5, 3)).toBe(2)
	})

	it('should multiply two numbers correctly', () => {
		expect(multiply(2, 3)).toBe(6)
	})

	it('should divide two numbers correctly', () => {
		const { quotient, remainder } = divide(7, 2)
		expect(quotient).toBe(3)
		expect(remainder).toBe(1)
	})

	it('should handle division by zero correctly', () => {
		const { quotient, remainder } = divide(7, 0)
		expect(quotient).toBe(0)
		expect(remainder).toBe(0)
	})
})
