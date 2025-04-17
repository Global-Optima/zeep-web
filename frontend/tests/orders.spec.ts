// tests/order-flows.spec.ts
import { expect, Page, test } from '@playwright/test'

const BASE_URL = 'https://dev.zeep.kz'

/**
 * Performs full Café‑user login, POS setup, and opens the kiosk menu.
 */
async function loginAndSetup(page: Page) {
	await page.goto(BASE_URL)

	// 1. Select café & employee, then log in
	await page.getByRole('button', { name: 'Кафе' }).click()
	await page.getByRole('combobox').filter({ hasText: 'Выберите кафе' }).click()
	await page.getByRole('option', { name: 'Центральное Кафе' }).click()
	await page.getByRole('combobox').filter({ hasText: 'Сотрудники' }).click()
	await page.getByRole('option', { name: 'Иван Иванов' }).click()
	await page.getByRole('textbox', { name: 'Введите пароль сотрудника' }).fill('Password1!')
	await page.getByRole('button', { name: 'Войти' }).click()

	// 2. Enter kiosk
	await page.getByRole('button', { name: 'Киоск' }).click()

	// 3. Initial POS‑terminal setup modal
	await expect(
		page.getByText('Перед началом работы убедитесь, что все пункты настроены'),
	).toBeVisible()
	await page.getByRole('button').click() // Confirm
	await page.getByRole('textbox', { name: 'Адрес POS-терминала' }).fill('Kaspi')
	await page.getByRole('textbox', { name: 'Имя интеграции' }).fill('kaspi')
	await page.getByRole('button', { name: 'Сохранить настройки и интегрировать' }).click()

	// 4. Re-enter kiosk landing
	await page.getByRole('button', { name: 'Киоск' }).click()
	await expect(page.getByText('Перейти в меню')).toBeVisible()
}

/**
 * Navigates from kiosk landing to the menu screen.
 */
async function openMenu(page: Page) {
	await page.locator('.bottom-0').first().click()
	await expect(page.getByText('Перейти в меню')).not.toBeVisible()
}

test.describe('Cafe Order E2E Flow', () => {
	test('ORDER‑01: Login & POS setup lands on menu', async ({ page }) => {
		await loginAndSetup(page)
		// Should show the “Перейти в меню” button
		await expect(page.getByText('Перейти в меню')).toBeVisible()
	})

	test('ORDER‑02: Search for “Мохито”', async ({ page }) => {
		await loginAndSetup(page)
		await openMenu(page)

		// Search “Мохито”
		await page.getByTestId('collapsed-search-trigger').click()
		const search = page.getByTestId('search-input')
		await search.fill('Мохито')
		await search.press('Enter')

		// Expect at least one product card containing “Мохито”
		const card = page.getByTestId('product-card').locator('div').filter({ hasText: 'Мохито' })
		await expect(card).toBeVisible()
	})

	test('ORDER‑03: Select M 350 мл variant and verify "700 ₸" add‑button', async ({ page }) => {
		await loginAndSetup(page)
		await openMenu(page)

		// Select the product
		await page.getByTestId('product-card').locator('div').filter({ hasText: 'Мохито' }).click()
		// Choose size M 350 мл
		await page.getByRole('button', { name: 'M 350 мл' }).click()
		// Verify add button shows exact price
		const addBtn = page
			.locator('div')
			.filter({ hasText: /^700 ₸$/ })
			.getByRole('button')

		await expect(addBtn).toBeVisible()
	})

	test('ORDER‑04: Place order with random name & Kaspi QR payment', async ({ page }) => {
		await loginAndSetup(page)
		await openMenu(page)

		// Select the product
		await page.getByTestId('product-card').locator('div').filter({ hasText: 'Мохито' }).click()
		// Choose size M 350 мл
		await page.getByRole('button', { name: 'M 350 мл' }).click()
		// Verify add button shows exact price
		await page
			.locator('div')
			.filter({ hasText: /^700 ₸$/ })
			.getByRole('button')
			.click()

		// Open order summary
		await page.getByRole('button', { name: '₸' }).click()
		// Randomize name
		await page.getByRole('button', { name: 'Случайное имя' }).click()
		// Proceed to payment
		await page.getByRole('button', { name: 'Оплатить' }).click()
		// Select Kaspi QR (2nd option)
		await page
			.locator('div')
			.filter({ hasText: /^Kaspi QR$/ })
			.nth(1)
			.click()

		// Expect final confirmation toast
		await expect(page.getByText('Спасибо за покупку! Заказ уже принят в работу')).toBeVisible()
	})

	test('ORDER‑05: Return to menu after order completion', async ({ page }) => {
		await loginAndSetup(page)
		await openMenu(page)

		// Select the product
		await page.getByTestId('product-card').locator('div').filter({ hasText: 'Мохито' }).click()
		// Choose size M 350 мл
		await page.getByRole('button', { name: 'M 350 мл' }).click()
		// Verify add button shows exact price
		await page
			.locator('div')
			.filter({ hasText: /^700 ₸$/ })
			.getByRole('button')
			.click()

		// Open order summary
		await page.getByRole('button', { name: '₸' }).click()
		// Randomize name
		await page.getByRole('button', { name: 'Случайное имя' }).click()
		// Proceed to payment
		await page.getByRole('button', { name: 'Оплатить' }).click()

		// Select Kaspi QR (2nd option)
		await page
			.locator('div')
			.filter({ hasText: /^Kaspi QR$/ })
			.nth(1)
			.click()
		await expect(page.getByText('Спасибо за покупку! Заказ уже принят в работу')).toBeVisible()

		// Click “Вернуться в меню”
		await page.getByRole('button', { name: 'Вернуться в меню' }).click()
		await expect(page.getByTestId('collapsed-search-trigger')).toBeVisible()
	})
})
