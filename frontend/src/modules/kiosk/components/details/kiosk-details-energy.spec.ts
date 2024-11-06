import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import KioskDetailsEnergy from './kiosk-details-energy.vue'

describe('KioskDetailsEnergy.vue', () => {
	const energy = {
		ccal: 250,
		proteins: 10,
		fats: 5,
		carbs: 20,
	}

	it('renders energy information correctly', () => {
		const wrapper = mount(KioskDetailsEnergy, {
			props: { energy },
		})

		// Check if each energy attribute is rendered with the correct value and label

		// Energy (ccal)
		const ccal = wrapper.find('[data-testid="energy-ccal"]')
		expect(ccal.exists()).toBe(true)
		expect(ccal.text()).toBe(`${energy.ccal} ккал`)

		// Protein
		const protein = wrapper.find('[data-testid="energy-protein"]')
		expect(protein.exists()).toBe(true)
		expect(protein.text()).toBe(`${energy.proteins} г`)

		// Fats
		const fats = wrapper.find('[data-testid="energy-fats"]')
		expect(fats.exists()).toBe(true)
		expect(fats.text()).toBe(`${energy.fats} г`)

		// Carbs
		const carbs = wrapper.find('[data-testid="energy-carbs"]')
		expect(carbs.exists()).toBe(true)
		expect(carbs.text()).toBe(`${energy.carbs} г`)
	})
})
