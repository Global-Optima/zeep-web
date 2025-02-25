package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Description  string
	Method       string
	URL          string
	Body         interface{}
	FormData     map[string]string
	Files        map[string][]*multipart.FileHeader
	Headers      map[string]string
	ExpectedCode int
	ExpectedBody interface{}
	AuthRole     data.EmployeeRole
}

func (env *TestEnvironment) RunTests(t *testing.T, testCases []TestCase) {
	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var req *http.Request
			var err error

			// Check if test case includes file uploads (multipart form-data)
			if len(tc.Files) > 0 || tc.FormData != nil {
				// Create multipart form-data request
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				// Attach form fields
				for key, val := range tc.FormData {
					_ = writer.WriteField(key, val)
				}

				// Attach file uploads
				for fieldName, fileHeaders := range tc.Files {
					for _, fileHeader := range fileHeaders {
						file, err := fileHeader.Open()
						if err != nil {
							t.Fatalf("failed to open test file: %v", err)
						}
						defer file.Close()

						part, err := writer.CreateFormFile(fieldName, fileHeader.Filename)
						if err != nil {
							t.Fatalf("failed to create form file: %v", err)
						}
						_, _ = io.Copy(part, file)
					}
				}

				_ = writer.Close()

				req, err = http.NewRequest(tc.Method, tc.URL, &buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				// Default JSON request
				body, _ := json.Marshal(tc.Body)
				req, err = http.NewRequest(tc.Method, tc.URL, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
			}

			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			if tc.AuthRole != "" {
				token := env.GetAuthToken(tc.AuthRole)
				if token != "" {
					req.Header.Set("Authorization", token)
				}
			}

			for key, value := range tc.Headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			env.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.ExpectedCode, w.Code)
			if tc.ExpectedBody != nil {
				expectedJSON, err := json.Marshal(tc.ExpectedBody)
				if err != nil {
					t.Fatalf("❌ Failed to marshal expected body: %v", err)
				}
				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			}
		})
	}
}

func (env *TestEnvironment) GetAuthToken(role data.EmployeeRole) string {
	// Check cache first
	if token, exists := env.Tokens[role.ToString()]; exists {
		return token
	}

	credentials := map[data.EmployeeRole]map[string]string{
		data.RoleAdmin:                  {"email": "jack@test.com", "password": "test"},
		data.RoleOwner:                  {"email": "ivy@test.com", "password": "test"},
		data.RoleStoreManager:           {"email": "alice@test.com", "password": "test"},
		data.RoleBarista:                {"email": "bob@test.com", "password": "test"},
		data.RoleWarehouseManager:       {"email": "emma@test.com", "password": "test"},
		data.RoleWarehouseEmployee:      {"email": "frank@test.com", "password": "test"},
		data.RoleFranchiseManager:       {"email": "henry@test.com", "password": "test"},
		data.RoleFranchiseOwner:         {"email": "ivy@test.com", "password": "test"},
		data.RoleRegionWarehouseManager: {"email": "grace@test.com", "password": "test"},
	}

	creds, exists := credentials[role]
	if !exists {
		log.Printf("⚠️ Warning: No credentials found for role %s", role)
		return ""
	}

	bodyBytes, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("❌ Failed to marshal login request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/test/auth/employees/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Printf("❌ Failed to authenticate as %s: Status %d", role, w.Code)
		return ""
	}

	var response struct {
		Data struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		log.Fatalf("❌ Failed to parse login response: %v", err)
	}

	env.Tokens[role.ToString()] = "Bearer " + response.Data.AccessToken
	return env.Tokens[role.ToString()]
}
