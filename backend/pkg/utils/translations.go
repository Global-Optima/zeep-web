package utils

import "gorm.io/gorm"

func WithLocalePreloads(locale string, relations ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, rel := range relations {
			db = db.Preload(rel, func(d *gorm.DB) *gorm.DB {
				return d.Where("language_code = ?", locale)
			})
		}
		return db
	}
}

func WithNameTranslations(locale string) func(*gorm.DB) *gorm.DB {
	return WithLocalePreloads(locale, "NameTranslation")
}

func WithAllTranslations(locale string) func(*gorm.DB) *gorm.DB {
	return WithLocalePreloads(locale, "NameTranslation", "DescriptionTranslation")
}
