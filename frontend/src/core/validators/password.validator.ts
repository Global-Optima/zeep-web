import * as z from 'zod'

// Regular expression for allowed characters (Latin letters, numbers, and special characters)
const LATIN_AND_SPECIAL_CHARS_REGEX = /^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*$/

// Reusable Zod schema for password validation
export const passwordValidationSchema = z
	.string()
	.min(8, 'Пароль должен содержать не менее 8 символов')
	.regex(
		LATIN_AND_SPECIAL_CHARS_REGEX,
		'Пароль должен содержать только латинские символы и разрешенные специальные символы',
	)
	.refine(
		password => /[A-Z]/.test(password),
		'Пароль должен содержать хотя бы одну заглавную букву',
	)
	.refine(password => /[a-z]/.test(password), 'Пароль должен содержать хотя бы одну строчную букву')
	.refine(password => /[0-9]/.test(password), 'Пароль должен содержать хотя бы одну цифру')
	.refine(
		password => /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
		'Пароль должен содержать хотя бы один специальный символ',
	)

export const easyPasswordValidationSchema = z
	.string()
	.min(2, 'Пароль должен содержать не менее 2 символов')
	.regex(
		LATIN_AND_SPECIAL_CHARS_REGEX,
		'Пароль должен содержать только латинские символы и разрешенные специальные символы',
	)
