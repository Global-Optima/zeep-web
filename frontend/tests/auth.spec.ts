// tests/auth.spec.ts
import { test, expect, Page } from '@playwright/test'

const BASE_URL = 'https://dev.zeep.kz/'

/**
 * Helper to select café, employee, fill password and click “Войти”
 */
async function doLogin(page: Page, cafe: string, employee: string, password: string) {
	await page.goto(BASE_URL)

	// Select café
	await page.getByRole('button', { name: 'Кафе' }).click()
	await page.getByRole('combobox').filter({ hasText: 'Выберите кафе' }).click()
	await page.getByRole('option', { name: cafe }).click()

	// Select employee
	await page.getByRole('combobox').filter({ hasText: 'Сотрудники' }).click()
	await page.getByRole('option', { name: employee }).click()

	// Enter password & submit
	await page.getByRole('textbox', { name: 'Введите пароль сотрудника' }).fill(password)
	await page.getByRole('button', { name: 'Войти' }).click()
}

test.describe('Authentication Scenarios', () => {
	test('AUTH‑01: Valid login succeeds', async ({ page }) => {
		await doLogin(page, 'Кафе на Углу', 'Сергей Павлов', 'Password1!')
		await expect(page.getByRole('heading', { name: /Добро пожаловать, Сергей/ })).toBeVisible()
	})

	test('AUTH‑02: Invalid password shows notification error', async ({ page }) => {
		await doLogin(page, 'Центральное Кафе', 'Иван Иванов', 'wrong-password')
		const notif = page.getByLabel('Notifications (F8)')
		await expect(
			notif.getByText(
				'Введённый вами адрес электронной почты или пароль неверны. Пожалуйста, попробуйте снова.',
			),
		).toBeVisible()
	})

	test('AUTH‑03: Disallowed characters show pattern‑validation error', async ({ page }) => {
		await page.goto(BASE_URL)
		// select café & employee
		await page.getByRole('button', { name: 'Кафе' }).click()
		await page.getByRole('combobox').filter({ hasText: 'Выберите кафе' }).click()
		await page.getByRole('option', { name: 'Центральное Кафе' }).click()
		await page.getByRole('combobox').filter({ hasText: 'Сотрудники' }).click()
		await page.getByRole('option', { name: 'Иван Иванов' }).click()

		// fill SQL‑injection style password
		await page.getByRole('textbox', { name: 'Введите пароль сотрудника' }).fill("' OR 1=1; --")

		// assert inline pattern message
		await expect(
			page.getByText(
				'Пароль должен содержать только латинские символы и разрешенные специальные символы',
			),
		).toBeVisible()
	})

	test('AUTH‑04: Too‑short password shows inline validation', async ({ page }) => {
		await page.goto(BASE_URL)
		// select café & employee
		await page.getByRole('button', { name: 'Кафе' }).click()
		await page.getByRole('combobox').filter({ hasText: 'Выберите кафе' }).click()
		await page.getByRole('option', { name: 'Кафе на Углу' }).click()
		await page.getByRole('combobox').filter({ hasText: 'Сотрудники' }).click()
		await page.getByRole('option', { name: 'Сергей Павлов' }).click()

		// fill too-short password
		await page.getByRole('textbox', { name: 'Введите пароль сотрудника' }).fill('1234567')

		// assert inline length message
		await expect(page.getByText('Пароль должен содержать не менее 8 символов')).toBeVisible()
	})
})
