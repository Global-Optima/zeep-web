import { expect, Page, test } from '@playwright/test'

const BASE_URL = 'https://dev.zeep.kz'

/**
 * Full Café‑user login + initial POS‑terminal configuration + open kiosk.
 */
async function loginCafe(page: Page, cafe: string, employee: string, password: string) {
	await page.goto(BASE_URL)

	// 1. Select café and employee, then log in
	await page.getByRole('button', { name: 'Кафе' }).click()
	await page.getByRole('combobox').filter({ hasText: 'Выберите кафе' }).click()
	await page.getByRole('option', { name: cafe }).click()
	await page.getByRole('combobox').filter({ hasText: 'Сотрудники' }).click()
	await page.getByRole('option', { name: employee }).click()
	await page.getByRole('textbox', { name: 'Введите пароль сотрудника' }).fill(password)
	await page.getByRole('button', { name: 'Войти' }).click()

	// 2. Click “Киоск” and handle initial setup modal
	await page.getByRole('button', { name: 'Киоск' }).click()
	await expect(
		page.getByText('Перед началом работы убедитесь, что все пункты настроены'),
	).toBeVisible()
	await page.getByText('Перед началом работы убедитесь, что все пункты настроены').click()
	await page.getByRole('button').click() // Confirm

	// 3. Fill POS‑terminal settings
	await page.getByRole('textbox', { name: 'Адрес POS-терминала' }).fill('kaspi.test.kz')
	await page.getByRole('textbox', { name: 'Имя интеграции' }).fill('Test')
	await page.getByRole('button', { name: 'Сохранить настройки и интегрировать' }).click()

	// 4. Re-enter kiosk
	await page.getByRole('button', { name: 'Киоск' }).click()
	await expect(page.getByText('Перейти в меню')).toBeVisible()
}

test.describe('Cafe Order Flows', () => {
	test.beforeEach(async ({ page }) => {
		await loginCafe(page, 'Центральное Кафе', 'Иван Иванов', 'Password1!')
		// Now on kiosk landing with “Перейти в меню”
		await page.getByText('Перейти в меню').click()
	})

	test('ORDER‑01: Search + add “Мохито”', async ({ page }) => {
		// Search for “Мохито”
		await page.getByTestId('collapsed-search-trigger').click()
		await page.getByTestId('search-input').fill('Мохито')
		await page.getByTestId('search-input').press('Enter')

		// Select and add the product
		await page.getByTestId('product-title').click()
		await page.getByText('Мохито').click()
		await page.getByRole('button').first().click()

		// Verify it appears in order summary
		await expect(page.getByText('Мохито')).toBeVisible()
	})

	test('ORDER‑02: Add “Клубничный смузи” with size & additives', async ({ page }) => {
		// Navigate to Smuzi category
		await page.getByRole('button', { name: 'Смузи' }).click()
		await page.getByText('Клубничный смузи1 500 ₸').click()
		await page.getByText('Клубничный смузи').click()

		// Choose size L
		await page.getByRole('button', { name: 'L 500 миллилитр' }).click()

		// Open additives panel and add 2× ice
		await page.getByRole('button', { name: '500 ₸' }).click()
		await page.getByTestId('additive-card').click()
		await page.getByText('Кубики льда50 ₸').click()
		await page.getByText('Кубики льда50 ₸').click()

		// Confirm total and add to cart
		const addBtn = page.getByRole('button', { name: '550 ₸' })
		await expect(addBtn).toBeVisible()
		await addBtn.click()

		// Verify in summary
		await expect(page.getByText('550 ₸')).toBeVisible()
	})

	test('ORDER‑03: Proceed to payment and back preserves order', async ({ page }) => {
		// Click “Оплатить”
		await page.getByRole('button', { name: 'Оплатить' }).click()
		await expect(page.getByRole('heading', { name: 'Выберите способ оплаты' })).toBeVisible()

		// Go back
		await page.getByRole('button', { name: 'Назад' }).click()
		await expect(page.getByRole('button', { name: '550 ₸' })).toBeVisible()
	})

	test('ORDER‑04: Complete via Kaspi QR and return to menu', async ({ page }) => {
		// Start payment
		await page.getByRole('button', { name: 'Оплатить' }).click()
		// Select Kaspi QR
		await page
			.locator('div')
			.filter({ hasText: /^Kaspi QR$/ })
			.nth(1)
			.click()
		await expect(page.getByRole('heading', { name: 'Ожидание оплаты' })).toBeVisible()
		await expect(page.getByRole('heading', { name: 'Заказ принят' })).toBeVisible()

		// Confirmation toast
		await expect(page.getByText('Спасибо за покупку! Заказ уже принят в работу')).toBeVisible()

		// Return to menu
		await page.getByRole('button', { name: 'Вернуться в меню' }).click()
		await expect(page.getByText('Перейти в меню')).toBeVisible()
	})
})
