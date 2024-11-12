package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Description  string
	Method       string
	URL          string
	Body         interface{}
	Headers      map[string]string
	ExpectedCode int
	ExpectedBody interface{}
}

func TestRunner(t *testing.T, router *gin.Engine, testCases []TestCase) {
	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var reqBody io.Reader

			if tc.Body != nil {
				bodyBytes, err := json.Marshal(tc.Body)
				if err != nil {
					t.Fatalf("Failed to marshal body: %v", err)
				}
				reqBody = bytes.NewReader(bodyBytes)
			}

			req := httptest.NewRequest(tc.Method, tc.URL, reqBody)
			req.Header.Set("Content-Type", "application/json")

			for key, value := range tc.Headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.ExpectedCode, w.Code)

			if tc.ExpectedBody != nil {
				expectedJSON, err := json.Marshal(tc.ExpectedBody)
				if err != nil {
					t.Fatalf("Failed to marshal expected body: %v", err)
				}
				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			}
		})
	}
}
