package testUtils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Name       string
	InputArgs  []interface{}
	Expected   interface{}
	ShouldFail bool
	SetupMock  func()
}

type TestFunc func(args ...interface{}) (interface{}, error)

func TestRunner(t *testing.T, testFunc TestFunc, testCases []TestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			if tc.SetupMock != nil {
				tc.SetupMock()
			}

			actual, err := testFunc(tc.InputArgs...)

			if tc.ShouldFail {
				assert.Error(t, err, "Expected test to fail, but it did not")
				assert.Nil(t, actual)
			} else {
				assert.NoError(t, err, "Expected test to succeed, but it failed")
				if actual == nil {
					actual = []interface{}{}
				}

				assert.Equal(t, tc.Expected, actual, "Unexpected result")
			}
		})
	}
}
