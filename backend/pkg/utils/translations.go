package utils

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

func WithLocalePreloads(locale data.LanguageCode, relations ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, rel := range relations {
			db = db.Preload(rel, func(d *gorm.DB) *gorm.DB {
				return d.Where("language_code = ?", locale)
			})
		}
		return db
	}
}

type LocalizedPreload struct {
	Relation  string
	Localized bool               // true = WHERE language_code = ?, false = just preload
	Nested    []LocalizedPreload // for recursive preloads
}

func ApplyLocalizedPreloads(db *gorm.DB, locale data.LanguageCode, preloads []LocalizedPreload) *gorm.DB {
	for _, p := range preloads {
		if len(p.Nested) == 0 {
			if p.Localized {
				db = db.Preload(p.Relation, func(d *gorm.DB) *gorm.DB {
					return d.Where("language_code = ?", locale)
				})
			} else {
				db = db.Preload(p.Relation)
			}
		} else {
			db = db.Preload(p.Relation, func(d *gorm.DB) *gorm.DB {
				return ApplyLocalizedPreloads(d, locale, p.Nested)
			})
		}
	}
	return db
}
