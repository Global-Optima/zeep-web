package data

type LanguageCode string

var (
	LanguageCodeEN LanguageCode = "en"
	LanguageCodeRU LanguageCode = "ru"
	LanguageCodeKK LanguageCode = "kk"
)

func (lc LanguageCode) IsValid() bool {
	switch lc {
	case LanguageCodeEN, LanguageCodeRU, LanguageCodeKK:
		return true
	default:
		return false
	}
}

func (lc LanguageCode) String() string {
	return string(lc)
}

type AppTranslations struct {
	BaseEntity
	TranslationID  uint         `gorm:"index;not null"`
	LanguageCode   LanguageCode `json:"language_code" gorm:"size:10;not null"`
	TranslatedText string       `json:"translated_text" gorm:"type:text"`
}
