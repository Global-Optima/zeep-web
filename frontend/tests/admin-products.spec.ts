import { test, expect, Page } from '@playwright/test'

const BASE_URL = 'https://dev.zeep.kz'

/**
 * Log in as administrator (Елена Соколова) and open the Products page.
 */
async function loginAsAdmin(page: Page) {
	await page.goto(BASE_URL)
	await page.getByRole('button', { name: 'Администратор' }).click()
	await page.getByRole('combobox').filter({ hasText: 'Учетные записи' }).click()
	await page.getByRole('option', { name: 'Елена Соколова' }).click()
	await page.getByRole('textbox', { name: 'Введите ваш пароль' }).fill('Password1!')
	await page.getByRole('button', { name: 'Войти' }).click()
	// verify login succeeded
	await expect(page.getByRole('heading', { name: /Добро пожаловать/ })).toBeVisible()

	// navigate to products
	await page.getByRole('button', { name: 'Продукты', exact: true }).click()
	await expect(page).toHaveURL(/\/admin\/products/)
}

test.describe('Admin Product Search – Vital Scenarios', () => {
	test.beforeEach(async ({ page }) => {
		await loginAsAdmin(page)
	})

	test('PS‑01: Exact‑match search returns only that product', async ({ page }) => {
		const search = page.getByRole('searchbox', { name: 'Поиск' })
		await search.fill('Карамельный латте')
		await search.press('Enter')

		// Only one row with exactly that name
		const cells = page.getByRole('cell', { name: 'Карамельный латте' })
		await expect(cells).toHaveCount(1)
	})

	test('PS‑02: Partial‑match search returns all containing items', async ({ page }) => {
		const search = page.getByRole('searchbox', { name: 'Поиск' })
		await search.fill('латте')
		await search.press('Enter')

		// Every visible product name should include "латте"
		const rows = page.getByRole('row').filter({ has: page.getByText(/латте/i) })
		await expect(rows).not.toHaveCount(0)
		for (const row of await rows.elementHandles()) {
			const text = await row.textContent()
			expect(text?.toLowerCase()).toContain('латте')
		}
	})

	test('PS‑03: Case‑insensitive search yields same count as lowercase', async ({ page }) => {
		const search = page.getByRole('searchbox', { name: 'Поиск' })

		// Uppercase
		await search.fill('КАРАМЕЛЬНЫЙ')
		await search.press('Enter')
		const upperCount = await page
			.getByRole('cell')
			.filter({ hasText: /карамельный латте/i })
			.count()

		// Reset by reloading (session persists)
		await page.reload()
		await expect(page).toHaveURL(/\/admin\/products/)

		// Lowercase
		await search.fill('карамельный')
		await search.press('Enter')
		const lowerCount = await page
			.getByRole('cell')
			.filter({ hasText: /карамельный латте/i })
			.count()

		expect(upperCount).toBe(lowerCount)
	})

	test('PS‑04: No‑results shows friendly message', async ({ page }) => {
		const search = page.getByRole('searchbox', { name: 'Поиск' })
		await search.fill('nonexistent123')
		await search.press('Enter')

		await expect(page.getByText('Продукты не найдены')).toBeVisible()
	})

	test('PS‑05: Special‑chars / SQL injection input is sanitized', async ({ page }) => {
		const search = page.getByRole('searchbox', { name: 'Поиск' })
		await search.fill("'; DROP TABLE products;--")
		await search.press('Enter')

		// No crash, just no results
		await expect(page.getByText('Продукты не найдены')).toBeVisible()
	})
})
