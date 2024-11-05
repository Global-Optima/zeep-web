import { formatPrice } from '@/core/utils/price.utils'
import type { Products } from '@/modules/products/models/product.model'
import { router } from '@/router'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import KioskHomeProductCard from './kiosk-home-product-card.vue'

describe('KioskHomeProductCard.vue', () => {
	const product: Products = {
		id: 1,
		title: 'Test Product',
		image: 'https://example.com/image.jpg',
		price: 1999.99,
		category: 'Кофе',
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
		expect(img.attributes('src')).toBe(product.image)
		expect(img.attributes('alt')).toBe('Product Image')

		const title = wrapper.find('[data-testid="product-title"]')
		expect(title.exists()).toBe(true)
		expect(title.text()).toBe(product.title)

		const priceText = formatPrice(product.price)
		const price = wrapper.find('[data-testid="product-price"]')
		expect(price.exists()).toBe(true)
		expect(price.text()).toBe(priceText)
	})
})
