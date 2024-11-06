import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import KioskHomeToolbar from './kiosk-home-toolbar.vue'

describe('KioskHomeToolbar.vue', () => {
	const categories = ['Coffee', 'Tea', 'Juice']

	it('renders the search input and buttons correctly', () => {
		const wrapper = mount(KioskHomeToolbar, {
			props: {
				categories,
				searchTerm: '',
				selectedCategory: '',
			},
		})

		const searchSection = wrapper.find('[data-testid="search-section"]')
		expect(searchSection.exists()).toBe(true)

		const searchInput = wrapper.find('[data-testid="search-input"]')
		expect(searchInput.exists()).toBe(true)
		expect(searchInput.attributes('placeholder')).toBe('Поиск')

		const categoryButtons = wrapper.findAll('[data-testid="category-button"]')
		expect(categoryButtons.length).toBe(categories.length)
		categoryButtons.forEach((button, index) => {
			expect(button.text()).toBe(categories[index])
		})
	})

	it('emits update:searchTerm when input value changes', async () => {
		const wrapper = mount(KioskHomeToolbar, {
			props: {
				categories,
				searchTerm: '',
				selectedCategory: '',
			},
		})

		const searchInput = wrapper.find('[data-testid="search-input"]')
		await searchInput.setValue('Espresso')

		expect(wrapper.emitted('update:searchTerm')).toBeTruthy()
		expect(wrapper.emitted('update:searchTerm')![0]).toEqual(['Espresso'])
	})

	it('emits update:searchTerm with an empty string when the clear button is clicked', async () => {
		const wrapper = mount(KioskHomeToolbar, {
			props: {
				categories,
				searchTerm: 'Espresso',
				selectedCategory: '',
			},
		})

		const clearButton = wrapper.find('[data-testid="clear-button"]')
		expect(clearButton.exists()).toBe(true)

		await clearButton.trigger('click')

		expect(wrapper.emitted('update:searchTerm')).toBeTruthy()
		expect(wrapper.emitted('update:searchTerm')![0]).toEqual([''])
	})

	it('emits update:category and clears search term when a category button is clicked', async () => {
		const wrapper = mount(KioskHomeToolbar, {
			props: {
				categories,
				searchTerm: '',
				selectedCategory: '',
			},
		})

		const categoryButtons = wrapper.findAll('[data-testid="category-button"]')
		expect(categoryButtons.length).toBe(categories.length)

		await categoryButtons[0].trigger('click')

		expect(wrapper.emitted('update:category')).toBeTruthy()
		expect(wrapper.emitted('update:category')![0]).toEqual(['Coffee'])

		expect(wrapper.emitted('update:searchTerm')).toBeTruthy()
		expect(wrapper.emitted('update:searchTerm')![0]).toEqual([''])
	})
})
