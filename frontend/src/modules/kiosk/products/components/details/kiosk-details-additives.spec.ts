import { formatPrice } from '@/core/utils/price.utils'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import type { Additive } from '../../models/product.model'
import KioskDetailsAdditives from './kiosk-details-additives.vue'

describe('KioskDetailsAdditives.vue', () => {
	const additive: Additive = {
		id: 1,
		name: 'Vanilla Syrup',
		imageUrl: 'https://example.com/vanilla.jpg',
		price: 500,
	}
	const selectedAdditives = [additive]

	it('renders the additive card correctly when not selected', () => {
		const wrapper = mount(KioskDetailsAdditives, {
			props: {
				additive,
				selectedAdditives: [],
				defaultAdditives: [],
			},
		})

		// Check if the additive card renders with default styles
		const card = wrapper.find('[data-testid="additive-card"]')
		expect(card.exists()).toBe(true)
		expect(card.classes()).toContain('bg-white')
		expect(card.classes()).toContain('border-transparent')

		// Check if image, name, and price are displayed
		const img = wrapper.find('[data-testid="additive-image"]')
		expect(img.attributes('src')).toBe(additive.imageUrl)

		const name = wrapper.find('[data-testid="additive-name"]')
		expect(name.text()).toBe(additive.name)

		const price = wrapper.find('[data-testid="additive-price"]')
		expect(price.text()).toBe(formatPrice(additive.price))
		expect(price.classes()).toContain('text-gray-400') // Price should be gray when not selected

		// Check button is not selected
		const button = wrapper.find('[data-testid="additive-button"]')
		expect(button.classes()).toContain('bg-gray-200')
		expect(wrapper.find('[data-testid="additive-selected-indicator"]').exists()).toBe(false)
	})

	it('renders the additive card correctly when selected', () => {
		const wrapper = mount(KioskDetailsAdditives, {
			props: {
				additive,
				selectedAdditives,
				defaultAdditives: [],
			},
		})

		// Check if the additive card renders with selected styles
		const card = wrapper.find('[data-testid="additive-card"]')
		expect(card.exists()).toBe(true)
		expect(card.classes()).toContain('bg-primary')
		expect(card.classes()).toContain('border-primary')

		// Check price is displayed in black when selected
		const price = wrapper.find('[data-testid="additive-price"]')
		expect(price.classes()).toContain('text-black')

		// Check button is selected
		const button = wrapper.find('[data-testid="additive-button"]')
		expect(button.classes()).toContain('bg-primary')

		// Check for the selected indicator within the button
		const selectedIndicator = wrapper.find('[data-testid="additive-selected-indicator"]')
		expect(selectedIndicator.exists()).toBe(true)
	})

	it('emits "click:additive" event with the additive data when clicked', async () => {
		const wrapper = mount(KioskDetailsAdditives, {
			props: {
				additive,
				selectedAdditives: [],
				defaultAdditives: [],
			},
		})

		// Trigger click on the additive card
		await wrapper.find('[data-testid="additive-card"]').trigger('click')

		// Verify that the event was emitted with the correct payload
		expect(wrapper.emitted('click:additive')).toBeTruthy()
		expect(wrapper.emitted('click:additive')![0]).toEqual([additive])
	})
})
