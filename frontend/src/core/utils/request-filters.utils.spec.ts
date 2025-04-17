// src/core/utils/buildRequestFilter.utils.spec.ts
import { describe, expect, it } from 'vitest'
import { buildRequestFilter } from './request-filters.utils'

describe('buildRequestFilter', () => {
	it('returns an empty object when input is undefined', () => {
		expect(buildRequestFilter(undefined)).toEqual({})
	})

	it('filters out null, undefined, and empty-string values at the top level', () => {
		const input = {
			a: '',
			b: 'hello',
			c: null,
			d: undefined,
			e: 0,
			f: false,
		}
		expect(buildRequestFilter(input)).toEqual({
			b: 'hello',
			e: 0,
			f: false,
		})
	})

	it('recursively filters nested objects but preserves non-empty children', () => {
		const input = {
			profile: {
				name: '',
				email: 'user@example.com',
				age: null,
			},
			status: 'active',
		}
		expect(buildRequestFilter(input)).toEqual({
			profile: { email: 'user@example.com' },
			status: 'active',
		})
	})

	it('retains a nested object even if all its fields are removed (resulting in `{}`)', () => {
		const input = {
			metadata: {
				foo: '',
				bar: '',
			},
			tag: undefined,
		}
		expect(buildRequestFilter(input)).toEqual({
			metadata: {},
		})
	})

	it('does not recurse/filter inside arrays (arrays are left intact)', () => {
		const input = {
			arr: [null, '', 'valid', 0],
			nested: {
				list: ['', null],
			},
		}
		expect(buildRequestFilter(input)).toEqual({
			arr: [null, '', 'valid', 0],
			nested: { list: ['', null] },
		})
	})
})
