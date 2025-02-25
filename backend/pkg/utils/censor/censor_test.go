package censor

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func TestMain(m *testing.M) {
	err := InitCensor()
	if err != nil {
		logrus.Fatal(err)
		panic("Failed to initialize censor: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestNormalizeTextFunc(t *testing.T) {
	tests := []struct {
		description    string
		input          string
		expectedOutput string
	}{
		{
			description:    "Underlining should be removed",
			input:          "SingleWordTESTCASe_1",
			expectedOutput: "singlewordtestcase1",
		},
		{
			description:    "Spaces and underlining should be removed",
			input:          "Single  Word TEST CASe_1",
			expectedOutput: "singlewordtestcase1",
		},
		{
			description:    "Multiple spaces and underlining should be removed",
			input:          "оу=;562 ва һөҢ,ғsflk__Er",
			expectedOutput: "оу=;562ваһөң,ғsflker",
		},
		{
			description:    "Valid regular input from frontend",
			input:          "Покупатель-895",
			expectedOutput: "покупатель895",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actualOutput := normalizeText(tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}

func TestValidateTextFunc(t *testing.T) {
	validator := GetCensorValidator()

	tests := []struct {
		description string
		str         string
		expectError bool
	}{
		{
			description: "Specialized symbols should fail",
			str:         "сadӘңҺһө!#25^&",
			expectError: true,
		},
		{
			description: "Allowed specialized symbols should pass",
			str:         "сadӘңҺһ _ - ө25",
			expectError: false,
		},
		{
			description: "Frontend auto-generated 4 numbered name should pass",
			str:         "Покупатель-8456",
			expectError: false,
		},
		{
			description: "Frontend auto-generated 3 numbered name should pass",
			str:         "покупатель895",
			expectError: false,
		},
		{
			description: "Frontend auto-generated 3 numbered name should pass",
			str:         "Покупатель-895",
			expectError: false,
		},
		{
			description: "Frontend auto-generated 2 numbered name should pass",
			str:         "Покупатель-89",
			expectError: false,
		},
		{
			description: "Frontend auto-generated 1 numbered name should pass",
			str:         "Покупатель-8",
			expectError: false,
		},
		{
			description: "Empty string with spaces should fail",
			str:         "   ",
			expectError: true,
		},
		{
			description: "English Bad word should fail",
			str:         "Fu_cker",
			expectError: true,
		},
		{
			description: "Russian bad word should fail",
			str:         "гандөн228",
			expectError: true,
		},
		{
			description: "Kazakh bad word should fail",
			str:         "коtak",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := validator.ValidateText(tt.str)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegexp(t *testing.T) {
	re := regexp.MustCompile(ALLOWED_SYMBOLS)
	tests := []struct {
		description string
		input       string
		regexp      *regexp.Regexp
		isAllowed   bool
	}{
		{
			description: "Valid input",
			input:       "ads_-__фывөқӘҢ",
			regexp:      re,
			isAllowed:   true,
		},
		{
			description: "Invalid input with special symbols",
			input:       "ads_-__фывө@$1_+=$қӘҢ",
			regexp:      re,
			isAllowed:   false,
		},
		{
			description: "Valid regular input from frontend",
			input:       "Покупатель-5486",
			regexp:      re,
			isAllowed:   true,
		},
		{
			description: "Valid regular input from frontend",
			input:       "Покупатель-895",
			regexp:      re,
			isAllowed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			matchRegexp := re.MatchString(tt.input)
			assert.Equal(t, tt.isAllowed, matchRegexp)
		})
	}
}
