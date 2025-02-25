package utils_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	testUtils "github.com/Global-Optima/zeep-web/backend/tests/unit/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Constants for test data
const (
	validPhone   = "+77011234567"
	invalidPhone = "12345"
	emptyPhone   = ""
	validEmail   = "test@example.com"
	invalidEmail = "invalid-email"
	emptyEmail   = ""
	multiAtEmail = "test@@example.com"
	countryKZ    = "KZ"
)

func TestValidateTime(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Time Format", InputArgs: []interface{}{"23:59"}, Expected: []interface{}{}, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Time Format", InputArgs: []interface{}{"25:00"}, Expected: []interface{}{}, ShouldFail: true, SetupMock: nil},
	}

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return nil, utils.ValidateTime(args[0].(string))
	}, testCases)
}

func TestIsValidPhone(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Phone", InputArgs: []interface{}{validPhone, countryKZ}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Phone", InputArgs: []interface{}{invalidPhone, countryKZ}, Expected: false, ShouldFail: false, SetupMock: nil},
		{Name: "Empty Phone", InputArgs: []interface{}{emptyPhone, countryKZ}, Expected: false, ShouldFail: false, SetupMock: nil},
	}

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidPhone(args[0].(string), args[1].(string)), nil
	}, testCases)
}

func TestIsValidEmail(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Email", InputArgs: []interface{}{validEmail}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Email", InputArgs: []interface{}{invalidEmail}, Expected: false, ShouldFail: false, SetupMock: nil},
		{Name: "Empty Email", InputArgs: []interface{}{emptyEmail}, Expected: false, ShouldFail: false, SetupMock: nil},
		{Name: "Multiple @ Email", InputArgs: []interface{}{multiAtEmail}, Expected: false, ShouldFail: false, SetupMock: nil},
	}

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidEmail(args[0].(string)), nil
	}, testCases)
}

func TestIsValidLatitudeLongitude(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Latitude", InputArgs: []interface{}{45.0}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Latitude", InputArgs: []interface{}{-100.0}, Expected: false, ShouldFail: false, SetupMock: nil},
		{Name: "Valid Longitude", InputArgs: []interface{}{90.0}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Longitude", InputArgs: []interface{}{190.0}, Expected: false, ShouldFail: false, SetupMock: nil},
	}

	// Run latitude tests
	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidLatitude(args[0].(float64)), nil
	}, testCases[:2])

	// Run longitude tests
	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidLongitude(args[0].(float64)), nil
	}, testCases[2:])
}

func TestParseParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})

	id, err := utils.ParseParam(c, "id")
	assert.NoError(t, err)
	assert.Equal(t, uint(123), id)
}
