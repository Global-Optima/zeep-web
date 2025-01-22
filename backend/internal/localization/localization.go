package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type LocalizedMessages struct {
	En string `json:"en"`
	Ru string `json:"ru"`
	Kk string `json:"kk"`
}

type Locale string

func (l Locale) ToString() string {
	return string(l)
}

func ToLocale(localeStr string) Locale {
	if strings.TrimSpace(localeStr) == "" {
		return DEFAULT_LOCALE
	}

	locale := Locale(strings.ToLower(localeStr))

	switch locale {
	case English, Russian, Kazakh:
		return locale
	}

	return DEFAULT_LOCALE
}

const (
	Kazakh  Locale = "kk"
	Russian Locale = "ru"
	English Locale = "en"

	DEFAULT_LOCALE = Russian
)

var localizers = make(map[Locale]*i18n.Localizer)
var defaultLocalizer *i18n.Localizer

func InitLocalizer() error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	dir := "./internal/localization/languages"

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			_, err := bundle.LoadMessageFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	for _, locale := range []Locale{Kazakh, Russian, English} {
		if locale == DEFAULT_LOCALE {
			defaultLocalizer = i18n.NewLocalizer(bundle, DEFAULT_LOCALE.ToString())
			localizers[locale] = defaultLocalizer
			continue
		}
		localizers[locale] = i18n.NewLocalizer(bundle, locale.ToString())
	}

	return nil
}

func Translate(messageID string, data map[string]interface{}) (*LocalizedMessages, error) {
	localizedMessages := &LocalizedMessages{}

	localizeCfg := &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	}

	enLocalizer, exists := localizers[English]
	if !exists {
		enLocalizer = localizers[DEFAULT_LOCALE]
	}
	ruLocalizer, exists := localizers[Russian]
	if !exists {
		ruLocalizer = localizers[DEFAULT_LOCALE]
	}
	kkLocalizer, exists := localizers[Kazakh]
	if !exists {
		kkLocalizer = localizers[DEFAULT_LOCALE]
	}

	var err error
	localizedMessages.En, err = enLocalizer.Localize(localizeCfg)
	if err != nil {
		return nil, err
	}
	localizedMessages.Ru, err = ruLocalizer.Localize(localizeCfg)
	if err != nil {
		return nil, err
	}
	localizedMessages.Kk, err = kkLocalizer.Localize(localizeCfg)
	if err != nil {
		return nil, err
	}

	return localizedMessages, nil
}

func FormTranslationKey(keys ...string) string {
	camelCaseKeys := make([]string, len(keys))
	for i, key := range keys {
		camelCaseKeys[i] = ToCamelCase(key)
	}
	return strings.Join(camelCaseKeys, ".")
}

func ToCamelCase(str string) string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = strings.ToLower(words[i])
		if i > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}

	return strings.Join(words, "")
}
