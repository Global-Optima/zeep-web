package localization

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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

var DefaultLocalizedErrorMessages = &LocalizedMessage{
	Ru: "Произошла непредвиденная ошибка. Пожалуйста, попробуйте позже.",
	En: "An unexpected error occurred. Please try again later.",
	Kk: "Күтпеген қате орын алды. Кейінірек қайтадан көріңіз.",
}

type LocalizedResponse struct {
	Message   *LocalizedMessage `json:"message"`
	Status    int               `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
}

func TranslateResponse(c *gin.Context, status int, componentName data.ComponentName) *LocalizedResponse {
	var localizedMessages *LocalizedMessage
	var err error

	if componentName != "" {
		localizedMessages, err = translateComponentResponse(status, componentName)
	}

	if err != nil || componentName == "" {
		localizedMessages, err = translateCommonResponse(status)
	}
	if err != nil {
		localizedMessages = DefaultLocalizedErrorMessages
	}

	return NewLocalizedResponse(c, localizedMessages, status)
}

func NewLocalizedResponse(c *gin.Context, localizedMessages *LocalizedMessage, status int) *LocalizedResponse {
	return &LocalizedResponse{
		Message:   localizedMessages,
		Status:    status,
		Timestamp: time.Now().UTC(),
		Path:      c.FullPath(),
	}
}

func translateComponentResponse(status int, componentName data.ComponentName) (*LocalizedMessage, error) {
	return Translate(FormTranslationKey(RESPONSES_KEY, strconv.Itoa(status), "-", componentName.ToString()),
		map[string]interface{}{})
}

func translateCommonResponse(status int) (*LocalizedMessage, error) {
	return Translate(FormTranslationKey(RESPONSES_KEY, strconv.Itoa(status)),
		map[string]interface{}{})
}
