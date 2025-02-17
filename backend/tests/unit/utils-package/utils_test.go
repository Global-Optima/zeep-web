package utils_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	testUtils "github.com/Global-Optima/zeep-web/backend/tests/unit/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateTime(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Time Format", InputArgs: []interface{}{"23:59"}, Expected: []interface{}([]interface{}{}), ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Time Format", InputArgs: []interface{}{"25:00"}, Expected: []interface{}([]interface{}{}), ShouldFail: true, SetupMock: nil},
	}

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return nil, utils.ValidateTime(args[0].(string))
	}, testCases)
}

func TestIsValidPhone(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Phone", InputArgs: []interface{}{"+77011234567", "KZ"}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Phone", InputArgs: []interface{}{"12345", "KZ"}, Expected: false, ShouldFail: false, SetupMock: nil},
	}

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidPhone(args[0].(string), args[1].(string)), nil
	}, testCases)
}

func TestIsValidEmail(t *testing.T) {
	testCases := []testUtils.TestCase{
		{Name: "Valid Email", InputArgs: []interface{}{"test@example.com"}, Expected: true, ShouldFail: false, SetupMock: nil},
		{Name: "Invalid Email", InputArgs: []interface{}{"invalid-email"}, Expected: false, ShouldFail: false, SetupMock: nil},
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

	testUtils.TestRunner(t, func(args ...interface{}) (interface{}, error) {
		return utils.IsValidLatitude(args[0].(float64)), nil
	}, testCases[:2])

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
