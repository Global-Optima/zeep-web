package functional

import (
	"fmt"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func uintVal(u uint) uint {
	return u
}

func orderStatusPtr(s string) *data.OrderStatus {
	status := data.OrderStatus(s)
	return &status
}

func TestOrderService_GetOrders_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetOrdersModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name            string
		filter          types.OrdersFilterQuery
		expectedCount   int
		expectedOrderID uint
	}{
		{
			name: "Search matches John",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr("John"),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 2,
		},
		{
			name: "Search matches Doe substring",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr("Doe"),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 2,
		},
		{
			name: "Search not found in Test Store 1",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr("Alice"),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 0,
		},
		{
			name: "Empty search returns all orders for Test Store 1",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr(""),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 2,
		},
		{
			name: "Filter by COMPLETED status",
			filter: types.OrdersFilterQuery{
				Status:     orderStatusPtr("COMPLETED"),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount:   1,
			expectedOrderID: 2,
		},
		{
			name: "Pagination: page 2 returns empty",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr(""),
				StoreID:    uintPtr(1),
				BaseFilter: tests.BaseFilterWithPagination(2, 10),
			},
			expectedCount: 0,
		},
		{
			name: "Filter by non-existent store",
			filter: types.OrdersFilterQuery{
				Search:     stringPtr("John"),
				StoreID:    uintPtr(999),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ordersFound, err := module.Service.GetOrders(tc.filter)
			assert.NoError(t, err, fmt.Sprintf("GetOrders returned an error in case %q", tc.name))
			assert.Len(t, ordersFound, tc.expectedCount, fmt.Sprintf("Expected %d orders in case %q", tc.expectedCount, tc.name))
		})
	}
}

func TestOrderService_GetAllBaristaOrders_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetOrdersModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	// nowUTC := time.Now().UTC()
	// startToday := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 12, 0, 0, 0, time.UTC)
	// updateQuery := fmt.Sprintf("UPDATE orders SET created_at = '%s' WHERE store_id = 1", startToday.Format(time.RFC3339))
	// if err := db.Exec(updateQuery).Error; err != nil {
	// 	t.Fatalf("Failed to update orders created_at: %v", err)
	// }

	testCases := []struct {
		name          string
		filter        types.OrdersTimeZoneFilter
		expectedCount int
	}{
		{
			name: "Valid orders using Asia/Qyzylorda",
			filter: types.OrdersTimeZoneFilter{
				StoreID:          uintPtr(1),
				TimeZoneLocation: stringPtr("Asia/Qyzylorda"),
			},

			expectedCount: 2,
		},
		{
			name: "Valid orders using TimeZoneOffset +0",
			filter: types.OrdersTimeZoneFilter{
				StoreID:        uintPtr(1),
				TimeZoneOffset: uintPtr(0),
			},
			expectedCount: 2,
		},
		{
			name: "No orders returned for America/New_York",
			filter: types.OrdersTimeZoneFilter{
				StoreID:          uintPtr(1),
				TimeZoneLocation: stringPtr("America/New_York"),
			},

			expectedCount: 0,
		},
		{
			name: "No orders for non-existent store",
			filter: types.OrdersTimeZoneFilter{
				StoreID: uintPtr(999),
			},
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ordersFound, err := module.Service.GetAllBaristaOrders(tc.filter)
			assert.NoError(t, err, fmt.Sprintf("GetAllBaristaOrders returned an error in case %q", tc.name))
			assert.Len(t, ordersFound, tc.expectedCount, fmt.Sprintf("Expected %d orders in case %q", tc.expectedCount, tc.name))
		})
	}
}
