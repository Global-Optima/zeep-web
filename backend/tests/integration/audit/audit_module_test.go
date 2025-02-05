package audit_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestAuditEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Fetch audit records", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of audits",
				Method:       http.MethodGet,
				URL:          "/audits",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a single audit record by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return an audit record",
				Method:       http.MethodGet,
				URL:          "/audits/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent audit record",
				Method:       http.MethodGet,
				URL:          "/audits/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Record an employee action", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should record an employee action",
				Method:      http.MethodPost,
				URL:         "/audits",
				Body: map[string]interface{}{
					"operationType": "CREATE",
					"componentName": "PRODUCT",
					"employeeId":    1,
					"ipAddress":     "192.168.1.1",
					"resourceUrl":   "/products",
					"method":        "POST",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Record multiple employee actions", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should record multiple employee actions",
				Method:      http.MethodPost,
				URL:         "/audits/bulk",
				Body: []map[string]interface{}{
					{
						"operationType": "UPDATE",
						"componentName": "PRODUCT",
						"employeeId":    1,
						"ipAddress":     "192.168.1.2",
						"resourceUrl":   "/products/1",
						"method":        "PUT",
					},
					{
						"operationType": "DELETE",
						"componentName": "PRODUCT",
						"employeeId":    1,
						"ipAddress":     "192.168.1.3",
						"resourceUrl":   "/products/2",
						"method":        "DELETE",
					},
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})
}
