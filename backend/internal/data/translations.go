package data

type LanguageCode string

var (
	LanguageCodeEN LanguageCode = "en"
	LanguageCodeRU LanguageCode = "ru"
	LanguageCodeKK LanguageCode = "kk"
)

type AppTranslations struct {
	BaseEntity
	TranslationID  uint         `gorm:"index;not null"`
	LanguageCode   LanguageCode `json:"language_code" gorm:"size:10;not null"`
	TranslatedText string       `json:"translated_text" gorm:"type:text"`
}
