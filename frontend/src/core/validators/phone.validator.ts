import { z } from 'zod'

// Reusable Zod schema for phone validation with separate error messages
export const phoneValidationSchema = z
	.string()
	.refine(value => value.startsWith('+'), "Номер телефона должен начинаться с символа '+'.")
	.refine(
		value => /^\+[1-9]\d{0,2}/.test(value),
		"Неверный код страны. Код должен состоять из 1–3 цифр и не может начинаться с '0'.",
	)
	.refine(
		value => /^\+[1-9]\d{0,2}[1-9]\d{9,}$/.test(value),
		'Номер телефона слишком короткий. Введите как минимум 10 цифр после кода страны.',
	)
	.refine(
		value => value.length <= 15,
		'Номер телефона слишком длинный. Максимальная длина — 15 символов.',
	)
