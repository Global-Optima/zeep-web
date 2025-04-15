package functional

import (
	"testing"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var container = tests.NewTestContainer()

// resetTestData resets the database by truncating tables and loading mock data.
func resetTestData(t *testing.T) *gorm.DB {
	db := container.GetDB()
	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}
	return db
}

// func stringPtr(s string) *string {
// 	return &s
// }

func uintPtr(u uint) *uint {
	return &u
}

// func orderStatusPtr(s string) *data.OrderStatus {
// 	status := data.OrderStatus(s)
// 	return &status
// }

// func TestOrderService_GetOrders_WithPreloadedData(t *testing.T) {
// 	_ = resetTestData(t)
// 	module := tests.GetOrdersModule()

// 	testCases := []struct {
// 		name            string
// 		filter          types.OrdersFilterQuery
// 		expectedCount   int
// 		expectedOrderID uint
// 	}{
// 		{
// 			name: "Search matches John",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr("John"),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount: 2,
// 		},
// 		{
// 			name: "Search matches Doe substring",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr("Doe"),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount: 2,
// 		},
// 		{
// 			name: "Search not found in Test Store 1",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr("Alice"),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount: 0,
// 		},
// 		{
// 			name: "Empty search returns all orders for Test Store 1",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr(""),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount: 2,
// 		},
// 		{
// 			name: "Filter by COMPLETED status",
// 			filter: types.OrdersFilterQuery{
// 				Status:     orderStatusPtr("COMPLETED"),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount:   1,
// 			expectedOrderID: 2,
// 		},
// 		{
// 			name: "Pagination: page 2 returns empty",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr(""),
// 				StoreID:    uintPtr(1),
// 				BaseFilter: tests.BaseFilterWithPagination(2, 10),
// 			},
// 			expectedCount: 0,
// 		},
// 		{
// 			name: "Filter by non-existent store",
// 			filter: types.OrdersFilterQuery{
// 				Search:     stringPtr("John"),
// 				StoreID:    uintPtr(999),
// 				BaseFilter: tests.BaseFilterWithPagination(1, 10),
// 			},
// 			expectedCount: 0,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ordersFound, err := module.Service.GetOrders(tc.filter)
// 			assert.NoError(t, err, fmt.Sprintf("GetOrders returned an error in case %q", tc.name))
// 			assert.Len(t, ordersFound, tc.expectedCount, fmt.Sprintf("Expected %d orders in case %q", tc.expectedCount, tc.name))
// 		})
// 	}
// }

// func TestOrderService_GetAllBaristaOrders_WithPreloadedData(t *testing.T) {
// 	container := tests.NewTestContainer()
// 	db := container.GetDB()
// 	module := tests.GetOrdersModule()

// 	if err := tests.TruncateAllTables(db); err != nil {
// 		t.Fatalf("Failed to truncate all tables: %v", err)
// 	}
// 	if err := tests.LoadTestData(db); err != nil {
// 		t.Fatalf("Failed to load test data: %v", err)
// 	}

// 	// Use the server's local time zone.
// 	localLocation := time.Now().Location()
// 	localNow := time.Now().In(localLocation)
// 	// Build a "today" timestamp at 00:01 local time.
// 	orderTimeLocal := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 1, 0, 0, localLocation)
// 	// Convert that to UTC for storing in the DB.
// 	orderCreatedUTC := orderTimeLocal.UTC()

// 	// Update orders for store 1 with the computed UTC timestamp.
// 	updateQuery := fmt.Sprintf("UPDATE orders SET created_at = '%s' WHERE store_id = 1", orderCreatedUTC.Format("2006-01-02 15:04:05"))
// 	if err := db.Exec(updateQuery).Error; err != nil {
// 		t.Fatalf("Failed to update orders created_at: %v", err)
// 	}

// 	// Pick an alternative real timezone.
// 	// If the local timezone is America/New_York, use Asia/Tokyo; otherwise, use America/New_York.
// 	var altTimezone string
// 	if localLocation.String() == "America/Los_Angeles" {
// 		altTimezone = "Asia/Tokyo"
// 	} else {
// 		altTimezone = "America/Los_Angeles"
// 	}

// 	testCases := []struct {
// 		name          string
// 		filter        types.OrdersTimeZoneFilter
// 		expectedCount int
// 	}{
// 		{
// 			name: "Valid orders using local timezone",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:          uintPtr(1),
// 				TimeZoneLocation: stringPtr(localLocation.String()),
// 			},
// 			// Using the local timezone, the order appears as created today.
// 			expectedCount: 2,
// 		},
// 		{
// 			name: "No orders returned for alternative timezone",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:          uintPtr(1),
// 				TimeZoneLocation: stringPtr(altTimezone),
// 			},
// 			// When converting orderCreatedUTC into the alternative timezone,
// 			// it should fall on the previous day (yesterday).
// 			expectedCount: 0,
// 		},
// 		{
// 			name: "No orders for non-existent store",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID: uintPtr(999),
// 			},
// 			expectedCount: 0,
// 		},
// 	}

// 	// Run test cases.
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ordersFound, err := module.Service.GetAllBaristaOrders(tc.filter)
// 			assert.NoError(t, err, fmt.Sprintf("GetAllBaristaOrders returned an error in case %q", tc.name))
// 			assert.Len(t, ordersFound, tc.expectedCount, fmt.Sprintf("Expected %d orders in case %q", tc.expectedCount, tc.name))
// 		})
// 	}
// }

// func TestOrderService_GetStatusesCount_WithPreloadedData(t *testing.T) {
// 	container := tests.NewTestContainer()
// 	db := container.GetDB()
// 	module := tests.GetOrdersModule()

// 	if err := tests.TruncateAllTables(db); err != nil {
// 		t.Fatalf("Failed to truncate all tables: %v", err)
// 	}
// 	if err := tests.LoadTestData(db); err != nil {
// 		t.Fatalf("Failed to load test data: %v", err)
// 	}

// 	// Use the server's local time zone.
// 	localLocation := time.Now().Location()
// 	localNow := time.Now().In(localLocation)
// 	// Build a "today" timestamp at 00:01 local time.
// 	orderTimeLocal := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 1, 0, 0, localLocation)
// 	orderCreatedUTC := orderTimeLocal.UTC()

// 	// Update orders for store 1 and store 2.
// 	updateQueryStore1 := fmt.Sprintf("UPDATE orders SET created_at = '%s' WHERE store_id = 1", orderCreatedUTC.Format("2006-01-02 15:04:05"))
// 	if err := db.Exec(updateQueryStore1).Error; err != nil {
// 		t.Fatalf("Failed to update orders created_at for store 1: %v", err)
// 	}
// 	updateQueryStore2 := fmt.Sprintf("UPDATE orders SET created_at = '%s' WHERE store_id = 2", orderCreatedUTC.Format("2006-01-02 15:04:05"))
// 	if err := db.Exec(updateQueryStore2).Error; err != nil {
// 		t.Fatalf("Failed to update orders created_at for store 2: %v", err)
// 	}

// 	// Pick an alternative real timezone as before.
// 	var altTimezone string
// 	if localLocation.String() == "America/New_York" {
// 		altTimezone = "Asia/Tokyo"
// 	} else {
// 		altTimezone = "America/New_York"
// 	}

// 	testCases := []struct {
// 		name           string
// 		filter         types.OrdersTimeZoneFilter
// 		expectedCounts types.OrderStatusesCountDTO
// 	}{
// 		{
// 			name: "Valid statuses using local timezone for store 1",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:          uintPtr(1),
// 				TimeZoneLocation: stringPtr(localLocation.String()),
// 			},
// 			// Based on the mock data for store 1 (one PENDING and one COMPLETED order).
// 			expectedCounts: types.OrderStatusesCountDTO{
// 				PENDING:     1,
// 				COMPLETED:   1,
// 				PREPARING:   0,
// 				IN_DELIVERY: 0,
// 				DELIVERED:   0,
// 				CANCELLED:   0,
// 				ALL:         2,
// 			},
// 		},
// 		{
// 			name: "No orders returned using alternative timezone for store 1",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:          uintPtr(1),
// 				TimeZoneLocation: stringPtr(altTimezone),
// 			},
// 			// In the alternative zone, the UTC time converts to yesterday.
// 			expectedCounts: types.OrderStatusesCountDTO{
// 				PENDING:     0,
// 				COMPLETED:   0,
// 				PREPARING:   0,
// 				IN_DELIVERY: 0,
// 				DELIVERED:   0,
// 				CANCELLED:   0,
// 				ALL:         0,
// 			},
// 		},
// 		{
// 			name: "No orders for non-existent store",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID: uintPtr(999),
// 			},
// 			expectedCounts: types.OrderStatusesCountDTO{
// 				PENDING:     0,
// 				COMPLETED:   0,
// 				PREPARING:   0,
// 				IN_DELIVERY: 0,
// 				DELIVERED:   0,
// 				CANCELLED:   0,
// 				ALL:         0,
// 			},
// 		},
// 	}

// 	// Execute test cases.
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			statuses, err := module.Service.GetStatusesCount(tc.filter)
// 			assert.NoError(t, err, fmt.Sprintf("GetStatusesCount returned an error in case %q", tc.name))
// 			assert.Equal(t, tc.expectedCounts.PENDING, statuses.PENDING, "Mismatch in PENDING count")
// 			assert.Equal(t, tc.expectedCounts.COMPLETED, statuses.COMPLETED, "Mismatch in COMPLETED count")
// 			assert.Equal(t, tc.expectedCounts.PREPARING, statuses.PREPARING, "Mismatch in PREPARING count")
// 			assert.Equal(t, tc.expectedCounts.IN_DELIVERY, statuses.IN_DELIVERY, "Mismatch in IN_DELIVERY count")
// 			assert.Equal(t, tc.expectedCounts.DELIVERED, statuses.DELIVERED, "Mismatch in DELIVERED count")
// 			assert.Equal(t, tc.expectedCounts.CANCELLED, statuses.CANCELLED, "Mismatch in CANCELLED count")
// 			assert.Equal(t, tc.expectedCounts.ALL, statuses.ALL, "Mismatch in ALL count")
// 		})
// 	}
// }

// Test GetOrderBySubOrder function.
func TestOrderService_GetOrderBySubOrder(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	order, err := module.Service.GetOrderBySubOrder(1)
	assert.NoError(t, err, "GetOrderBySubOrder should not error for an existing suborder")
	assert.NotNil(t, order, "Returned order should not be nil")
}

// Test GetOrderById function.
func TestOrderService_GetOrderById(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	orderDTO, err := module.Service.GetOrderById(1)
	assert.NoError(t, err, "GetOrderById should not error for an existing order")
	assert.Equal(t, uint(1), orderDTO.ID, "Order ID should be 1")
}

// Test GetOrderDetails function.
func TestOrderService_GetOrderDetails(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	details, err := module.Service.GetOrderDetails(1, nil)
	assert.NoError(t, err, "GetOrderDetails should not error for an existing order")
	assert.NotNil(t, details, "Expected order details for an existing order")
}

// Test ExportOrders function.
func TestOrderService_ExportOrders(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	now := time.Now()
	// Default date range: 72 hours ago to 24 hours ahead.
	defaultStart := now.Add(-72 * time.Hour)
	defaultEnd := now.Add(24 * time.Hour)

	// Same day: truncate current time to the day.
	sameDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Yesterday and today range.
	yesterday := sameDay.Add(-24 * time.Hour)
	today := sameDay

	// Define table-driven test cases.
	testCases := []struct {
		description   string
		filter        types.OrdersExportFilterQuery
		expectedCount int
	}{
		{
			description: "Basic export with start and end date (expect 2 orders in store 1)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &defaultStart,
				EndDate:   &defaultEnd,
				StoreID:   uintPtr(1),
			},
			expectedCount: 2,
		},
		{
			description: "Export with only startDate specified (expect 2 orders)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &defaultStart,
				// EndDate omitted
				StoreID: uintPtr(1),
			},
			expectedCount: 2,
		},
		{
			description: "Export with only endDate specified (expect 2 orders)",
			filter: types.OrdersExportFilterQuery{
				// StartDate omitted
				EndDate: &defaultEnd,
				StoreID: uintPtr(1),
			},
			expectedCount: 2,
		},
		{
			description: "Export with different language specified (e.g., 'kk'; expect 0 orders if none match)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &defaultStart,
				EndDate:   &defaultEnd,
				StoreID:   uintPtr(1),
				Language:  "kk",
			},
			expectedCount: 2,
		},
		{
			description: "Export when no orders exist (use non-existent store id)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &defaultStart,
				EndDate:   &defaultEnd,
				StoreID:   uintPtr(999), // assuming store 999 has no orders
			},
			expectedCount: 0,
		},
		{
			description: "Export orders when orders are status-changing (simulate transient state; expect 2 orders)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &defaultStart,
				EndDate:   &defaultEnd,
				StoreID:   uintPtr(1),
			},
			expectedCount: 2,
		},
		{
			description: "Export with same day for start and end (narrow range; expect 0 orders if no order falls exactly on that day)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &sameDay,
				EndDate:   &sameDay,
				StoreID:   uintPtr(1),
			},
			expectedCount: 0,
		},
		{
			description: "Export from yesterday to today (expect 1 order)",
			filter: types.OrdersExportFilterQuery{
				StartDate: &yesterday,
				EndDate:   &today,
				StoreID:   uintPtr(1),
			},
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		exports, err := module.Service.ExportOrders(&tc.filter)
		assert.NoError(t, err, "ExportOrders should not error: "+tc.description)
		assert.Equal(t, tc.expectedCount, len(exports), "Unexpected number of orders for test case: "+tc.description)
	}
}

// Test GetSubOrders function.
func TestOrderService_GetSubOrders(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	suborders, err := module.Service.GetSubOrders(1)
	assert.NoError(t, err, "GetSubOrders should not error for an existing order")
	// For order 1, our test data should have 1 suborder.
	assert.Equal(t, 1, len(suborders), "Expected 1 suborder for order 1")
}

func TestOrderService_CreateOrder_Combined(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	err := censor.InitializeCensorForTests()
	if err != nil {
		t.Fatalf("Failed to initialize censor for tests: %v", err)
	}

	testCases := []struct {
		name              string
		dto               types.CreateOrderDTO
		storeID           uint
		expectedError     bool
		expectedErrSubstr string
		expectedTotal     float64
	}{
		{
			name: "Successful order creation - John Doe, quantity 2",
			dto: types.CreateOrderDTO{
				CustomerName: "John",
				Suborders: []types.CreateSubOrderDTO{
					{
						StoreProductSizeID: 1, // product size with store_price 2.75
						Quantity:           2,
						StoreAdditivesIDs:  []uint{1}, // additive with store_price 0.55
					},
				},
				StoreID: 1,
			},
			storeID:       1,
			expectedError: false,
			// Expected total = 2 * (2.75 + 0(0.55 not counted since it is a default additive)) = 5.50
			expectedTotal: 5.5,
		},
		{
			name: "Failure due to empty suborders",
			dto: types.CreateOrderDTO{
				CustomerName: "John",
				Suborders:    []types.CreateSubOrderDTO{},
				StoreID:      1,
			},
			storeID:           1,
			expectedError:     true,
			expectedErrSubstr: "order can not be empty",
		},
		{
			name: "Failure due to censored customer name",
			dto: types.CreateOrderDTO{
				CustomerName: "fuck", // Assume this is rejected by the censor validator.
				Suborders: []types.CreateSubOrderDTO{
					{
						StoreProductSizeID: 1,
						Quantity:           1,
						StoreAdditivesIDs:  []uint{1},
					},
				},
				StoreID: 1,
			},
			storeID:           1,
			expectedError:     true,
			expectedErrSubstr: "name", // Expect an error message indicating censorship.
		},
		{
			name: "Failure due to invalid product size",
			dto: types.CreateOrderDTO{
				CustomerName: "John",
				Suborders: []types.CreateSubOrderDTO{
					{
						StoreProductSizeID: 9999, // non-existent ID
						Quantity:           1,
						StoreAdditivesIDs:  []uint{1},
					},
				},
				StoreID: 1,
			},
			storeID:           1,
			expectedError:     true,
			expectedErrSubstr: "not found",
		},
	}

	// Run each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := module.Service.CreateOrder(tc.storeID, &tc.dto)
			if tc.expectedError {
				assert.Error(t, err, "Expected error in case %q", tc.name)
				if tc.expectedErrSubstr != "" {
					assert.Contains(t, err.Error(), tc.expectedErrSubstr, "Error message should contain %q", tc.expectedErrSubstr)
				}
				assert.Nil(t, order, "Order should be nil when creation fails")
			} else {
				assert.NoError(t, err, "Unexpected error in case %q", tc.name)
				assert.NotNil(t, order, "Returned order should not be nil in case %q", tc.name)
				// Verify total and status.
				assert.InDelta(t, tc.expectedTotal, order.Total, 0.01, "Order total is incorrect")
				// assert.Equal(t, data.OrderStatusPending, order.Status, "New order should have PENDING status")
				// or
				assert.Equal(t, data.OrderStatusWaitingForPayment, order.Status, "New order should have WAITING_FOR_PAYMENT status")
			}
		})
	}
}

func insertTestOrderWithTwoSuborders(t *testing.T, db *gorm.DB) (orderID uint, suborderIDs []uint) {
	t.Helper()

	// We assume that the repository returns these fixed prices:
	const productSizePrice = 2.75
	const additivePrice = 0.55

	// Create the order.
	order := data.Order{
		CustomerID:      uintPtr(1),
		CustomerName:    "Test Order",
		StoreEmployeeID: uintPtr(1),
		StoreID:         1,
		Status:          data.OrderStatusPending,
		DisplayNumber:   1,
	}

	// Create two suborders (each one unit)
	suborder := func() data.Suborder {
		return data.Suborder{
			StoreProductSizeID: 1,
			Price:              productSizePrice + additivePrice,
			Status:             data.SubOrderStatusPending,
			SuborderAdditives: []data.SuborderAdditive{
				{
					StoreAdditiveID: 1,
					Price:           additivePrice,
				},
			},
		}
	}
	order.Suborders = []data.Suborder{suborder(), suborder()}
	order.Total = (productSizePrice + additivePrice) * 2

	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("Failed to insert test order: %v", err)
	}

	// Retrieve suborder IDs for later reference.
	var subs []data.Suborder
	if err := db.Where("order_id = ?", order.ID).Find(&subs).Error; err != nil {
		t.Fatalf("Failed to retrieve suborders: %v", err)
	}
	for _, s := range subs {
		suborderIDs = append(suborderIDs, s.ID)
	}
	return order.ID, suborderIDs
}

func TestOrderService_AdvanceSubOrderStatus(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	// Insert an order with two suborders.
	orderID, suborderIDs := insertTestOrderWithTwoSuborders(t, db)
	assert.Len(t, suborderIDs, 2, "Expected two suborders inserted")

	// Advance the second suborder from PENDING to PREPARING.
	dto, err := module.Service.SetNextSubOrderStatus(suborderIDs[1], nil)
	assert.NoError(t, err, "Advancing suborder status should succeed")
	assert.Equal(t, data.SubOrderStatusPreparing, dto.Status, "Suborder should transition to PREPARING")

	// Advance the same suborder from PREPARING to COMPLETED.
	dto, err = module.Service.SetNextSubOrderStatus(suborderIDs[1], nil)
	assert.NoError(t, err, "Advancing suborder status again should succeed")
	assert.Equal(t, data.SubOrderStatusCompleted, dto.Status, "Suborder should transition to COMPLETED")

	dto, err = module.Service.SetNextSubOrderStatus(suborderIDs[0], nil)
	assert.NoError(t, err, "Advancing suborder status should succeed")
	assert.Equal(t, data.SubOrderStatusPreparing, dto.Status, "Suborder should transition to PREPARING")

	// Advance the same suborder from PREPARING to COMPLETED.
	dto, err = module.Service.SetNextSubOrderStatus(suborderIDs[0], nil)
	assert.NoError(t, err, "Advancing suborder status again should succeed")
	assert.Equal(t, data.SubOrderStatusCompleted, dto.Status, "Suborder should transition to COMPLETED")

	order, err := module.Service.GetOrderById(orderID)
	assert.NoError(t, err)
	assert.Equal(t, data.OrderStatusCompleted, order.Status, "Order status should be COMPLETED when all suborders are complete")
}
