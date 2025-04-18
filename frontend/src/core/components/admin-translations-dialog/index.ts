export enum TranslationsLanguage {
	EN = 'en',
	KK = 'kk',
	RU = 'ru',
}

export interface TranslationFieldLocale {
	field: string
	label: string
	locales: {
		[key in TranslationsLanguage]?: string
	}
}
