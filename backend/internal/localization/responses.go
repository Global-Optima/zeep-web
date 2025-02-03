package localization

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"strconv"
)

type LocalizedErrorInterface interface {
	GetEn() string
	GetRu() string
	GetKk() string
}

const (
	RESPONSES_KEY      = "responses"
	COMPONENT_NAME_KEY = "Component"
	NO_COMPONENT       = ""
)

var DefaultLocalizedErrorMessages = &LocalizedMessages{
	Ru: "Произошла непредвиденная ошибка. Пожалуйста, попробуйте позже.",
	En: "An unexpected error occurred. Please try again later.",
	Kk: "Күтпеген қате орын алды. Кейінірек қайтадан көріңіз.",
}

// TODO replace with string with data.ComponentName
func TranslateResponse(status int, componentName string) *LocalizedMessages {
	localizedMessages := DefaultLocalizedErrorMessages
	var err error

	if componentName != "" {
		localizedMessages, err = translateComponentResponse(status, data.ComponentName(componentName))
	}

	if err != nil || componentName == "" {
		localizedMessages, err = translateCommonResponse(status)
	}

	if err != nil {
		return DefaultLocalizedErrorMessages
	}
	return localizedMessages
}

func translateComponentResponse(status int, componentName data.ComponentName) (*LocalizedMessages, error) {
	return Translate(FormTranslationKey(RESPONSES_KEY, strconv.Itoa(status), "-", componentName.ToString()),
		map[string]interface{}{})
}

func translateCommonResponse(status int) (*LocalizedMessages, error) {
	return Translate(FormTranslationKey(RESPONSES_KEY, strconv.Itoa(status)),
		map[string]interface{}{})
}
