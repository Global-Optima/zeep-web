package localization

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	path := "languages"
	err := InitLocalizer(&path)
	if err != nil {
		panic("Failed to initialize localizer: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestResponseKeyBuilding(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		component   data.ComponentName
		optionalKey []string
		expectedKey string
	}{
		{
			name:        "200 StoreEmployee update",
			status:      200,
			component:   data.StoreEmployeeComponent,
			optionalKey: []string{"update"},
			expectedKey: "responses.200-storeEmployee-update",
		},
		{
			name:        "201 StoreEmployee",
			status:      201,
			component:   data.StoreEmployeeComponent,
			optionalKey: []string{},
			expectedKey: "responses.201-storeEmployee",
		},
		{
			name:        "400 StoreStock onlyOneRequestPerDay",
			status:      400,
			component:   data.StoreStockComponent,
			optionalKey: []string{"ONLY_ONE_REQUEST_PER_DAY"},
			expectedKey: "responses.400-storeStock-onlyOneRequestPerDay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseKey := NewResponseKey(tt.status, tt.component, tt.optionalKey...)
			actualKey := FormResponseTranslationKey(responseKey)
			assert.Equal(t, tt.expectedKey, actualKey)
		})
	}
}

func TestTranslateComponentResponse(t *testing.T) {
	tests := []struct {
		name             string
		responseKey      *ResponseKey
		expectedMessages *LocalizedMessage
	}{
		{
			name: "StoreEmployee Update Success",
			responseKey: NewResponseKey(
				200,
				data.StoreEmployeeComponent,
				data.UpdateOperation.ToString(),
			),
			expectedMessages: &LocalizedMessage{
				En: "Store employee successfully updated.",
				Ru: "Сотрудник магазина успешно обновлен.",
				Kk: "Дүкен қызметкері сәтті жаңартылды.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizedMessage, err := TranslateComponentResponse(tt.responseKey)

			assert.NoError(t, err)
			assert.NotNil(t, localizedMessage)

			assert.Equal(t, tt.expectedMessages, localizedMessage)
		})
	}
}

func TestTranslateCommonResponse(t *testing.T) {
	tests := []struct {
		name             string
		status           int
		expectedMessages *LocalizedMessage
	}{
		{
			name:   "BadRequest Error",
			status: 400,
			expectedMessages: &LocalizedMessage{
				En: "Bad request. Please check your input and try again.",
				Ru: "Неверный запрос. Пожалуйста, проверьте ввод и попробуйте снова.",
				Kk: "Қате сұраныс. Параметрлерді тексеріп, қайта көріңіз.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizedMessage, err := TranslateCommonResponse(tt.status)

			assert.NoError(t, err)
			assert.NotNil(t, localizedMessage)

			assert.Equal(t, tt.expectedMessages, localizedMessage)
		})
	}
}
