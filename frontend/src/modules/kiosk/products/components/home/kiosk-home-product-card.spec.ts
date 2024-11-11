import { formatPrice } from '@/core/utils/price.utils'
import { router } from '@/router'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import type { StoreProducts } from '../../models/product.model'
import KioskHomeProductCard from './kiosk-home-product-card.vue'

describe('KioskHomeProductCard.vue', () => {
	const product: StoreProducts = {
		id: 1,
		name: 'Test Product',
		imageUrl: 'https://example.com/image.jpg',
		basePrice: 1999.99,
		category: 'Кофе',
		description: 'Test Product Description',
	}

	it('renders the product information correctly', () => {
		const wrapper = mount(KioskHomeProductCard, {
			props: { product },
			global: {
				plugins: [router],
			},
		})

		const img = wrapper.find('[data-testid="product-image"]')
		expect(img.exists()).toBe(true)
		expect(img.attributes('src')).toBe(product.imageUrl)
		expect(img.attributes('alt')).toBe('Product Image')

		const title = wrapper.find('[data-testid="product-title"]')
		expect(title.exists()).toBe(true)
		expect(title.text()).toBe(product.name)

		const priceText = formatPrice(product.basePrice)
		const price = wrapper.find('[data-testid="product-price"]')
		expect(price.exists()).toBe(true)
		expect(price.text()).toBe(priceText)
	})
})
