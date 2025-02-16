package localization

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LocalizedErrorInterface interface {
	GetEn() string
	GetRu() string
	GetKk() string
}

const (
	RESPONSES_KEY = "responses"
)

var (
	ErrMessageBindingJSON  = NewResponseKey(http.StatusBadRequest, "json")
	ErrMessageBindingQuery = NewResponseKey(http.StatusBadRequest, "query")
	ErrMessageGettingImage = NewResponseKey(http.StatusInternalServerError, "image", "upload")
	ErrMessageGettingVideo = NewResponseKey(http.StatusInternalServerError, "video", "upload")
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

type ResponseKey struct {
	Status        int
	ComponentName data.ComponentName
	OptionalKeys  []string
}

func NewResponseKey(status int, componentName data.ComponentName, optionalKeys ...string) *ResponseKey {
	return &ResponseKey{
		Status:        status,
		ComponentName: componentName,
		OptionalKeys:  optionalKeys,
	}
}

func translateResponseWithStatus(c *gin.Context, status int) *LocalizedResponse {
	localizedMessages, err := translateCommonResponse(status)
	if err != nil {
		localizedMessages = DefaultLocalizedErrorMessages
	}
	return NewLocalizedResponse(c, localizedMessages, status)
}

func translateResponseWithKey(c *gin.Context, responseKey *ResponseKey) *LocalizedResponse {
	var localizedMessages *LocalizedMessage
	var err error

	if responseKey.ComponentName != "" {
		localizedMessages, err = translateComponentResponse(responseKey)
	}

	if err != nil || responseKey.ComponentName == "" {
		localizedMessages, err = translateCommonResponse(responseKey.Status)
	}
	if err != nil {
		localizedMessages = DefaultLocalizedErrorMessages
	}

	return NewLocalizedResponse(c, localizedMessages, responseKey.Status)
}

func NewLocalizedResponse(c *gin.Context, localizedMessages *LocalizedMessage, status int) *LocalizedResponse {
	return &LocalizedResponse{
		Message:   localizedMessages,
		Status:    status,
		Timestamp: time.Now().UTC(),
		Path:      c.FullPath(),
	}
}

func translateComponentResponse(responseKey *ResponseKey) (*LocalizedMessage, error) {
	return Translate(formResponseTranslationKey(responseKey),
		map[string]interface{}{})
}

func translateCommonResponse(status int) (*LocalizedMessage, error) {
	return Translate(FormTranslationKey(RESPONSES_KEY, strconv.Itoa(status)),
		map[string]interface{}{})
}

func formResponseTranslationKey(responseKey *ResponseKey) string {
	key := FormTranslationKey(RESPONSES_KEY, strconv.Itoa(responseKey.Status))
	key += "-" + ToCamelCase(responseKey.ComponentName.ToString())
	for _, optionalKey := range responseKey.OptionalKeys {
		key += "-" + optionalKey
	}
	return key
}

func SendLocalizedResponseWithKey(c *gin.Context, responseKey *ResponseKey) {
	utils.SendResponseWithStatus(c, translateResponseWithKey(c, responseKey), responseKey.Status)
}

func SendLocalizedResponseWithStatus(c *gin.Context, status int) {
	utils.SendResponseWithStatus(c, translateResponseWithStatus(c, status), status)
}
