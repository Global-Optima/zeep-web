package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

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

func Translate(locale Locale, messageID string, data map[string]interface{}) (string, error) {
	localizer, exists := localizers[locale]
	if !exists {
		localizer = localizers[DEFAULT_LOCALE]
	}

	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}
