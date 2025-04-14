import en from '../src/core/locales/en.json'
import kk from '../src/core/locales/kk.json'
import ru from '../src/core/locales/ru.json'

/**
 * Recursively flattens the keys of a JSON object into an array of strings in "dot" notation.
 *
 * @param obj - The JSON object from which to gather all keys
 * @param prefix - Used internally for recursion; do not pass a value on initial call
 * @returns - Array of strings (keys in dot notation, e.g., "KIOSK.CART.TITLE")
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getKeys(obj: Record<string, any>, prefix = ''): string[] {
	let keys: string[] = []

	for (const key in obj) {
		if (typeof obj[key] === 'object' && obj[key] !== null) {
			// Recursively gather subkeys
			keys = keys.concat(getKeys(obj[key], `${prefix}${key}.`))
		} else {
			// Final leaf node, push the completed key
			keys.push(`${prefix}${key}`)
		}
	}

	return keys
}

const localeKeys = {
	en: getKeys(en),
	ru: getKeys(ru),
	kk: getKeys(kk),
}

interface LocaleDifferences {
	[locale: string]: {
		missingKeys: string[]
		extraKeys: string[]
	}
}

function compareKeysAllLocales(localeKeys: Record<string, string[]>): LocaleDifferences {
	// Gather all unique keys from all locales
	const allKeys = new Set(Object.values(localeKeys).flat())

	const differences: LocaleDifferences = {}

	// For each locale, find missing and extra keys
	for (const [locale, keys] of Object.entries(localeKeys)) {
		const localeKeySet = new Set(keys)

		// Missing keys: exist in `allKeys` but not in the locale
		const missingKeys = [...allKeys].filter(k => !localeKeySet.has(k))

		// Extra keys: exist in the locale but not in `allKeys`
		// (unusual but might happen if some file has completely new keys no one else has)
		const extraKeys = keys.filter(k => !allKeys.has(k))

		if (missingKeys.length || extraKeys.length) {
			differences[locale] = { missingKeys, extraKeys }
		}
	}

	return differences
}

const differences = compareKeysAllLocales(localeKeys)

if (Object.keys(differences).length > 0) {
	console.error('❌ The following locale issues were found:\n')
	for (const [locale, { missingKeys, extraKeys }] of Object.entries(differences)) {
		if (missingKeys.length > 0) {
			console.error(
				`Locale [${locale}] is missing these keys:\n  • ${missingKeys.join('\n  • ')}\n`,
			)
		}
		if (extraKeys.length > 0) {
			console.error(
				`Locale [${locale}] has extra (unused) keys:\n  • ${extraKeys.join('\n  • ')}\n`,
			)
		}
	}
	process.exit(1)
} else {
	console.log('✅ All locales have the same set of keys.')
}
