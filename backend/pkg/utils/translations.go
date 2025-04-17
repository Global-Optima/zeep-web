package utils

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

func FirstTranslation(tr []data.AppTranslations) string {
	if len(tr) > 0 {
		return tr[0].TranslatedText
	}
	return ""
}

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
	Localized bool
	Model     interface{} // scope WHERE to correct table
	Nested    []LocalizedPreload
}

func ApplyLocalizedPreloads(db *gorm.DB, locale data.LanguageCode,
	preloads []LocalizedPreload) *gorm.DB {
	for _, p := range preloads {
		if len(p.Nested) == 0 {
			if p.Localized {
				if p.Model == nil {
					panic("missing Model for localized preload: " + p.Relation)
				}
				db = db.Preload(p.Relation, func(d *gorm.DB) *gorm.DB {
					return d.Model(p.Model).
						Where("language_code = ?", locale)
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

func Translation(rel string) LocalizedPreload {
	return LocalizedPreload{
		Relation:  rel,
		Localized: true,
		Model:     &data.AppTranslations{},
	}
}
