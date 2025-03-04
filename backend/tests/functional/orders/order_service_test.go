package functional

import (
	"fmt"
	"testing"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
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

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func orderStatusPtr(s string) *data.OrderStatus {
	status := data.OrderStatus(s)
	return &status
}

func TestOrderService_GetOrders_WithPreloadedData(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

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
// 			name: "Valid orders using TimeZoneOffset +0 (UTC)",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:        uintPtr(1),
// 				TimeZoneOffset: uintPtr(0), // UTC
// 			},
// 			// In UTC the order remains 00:01.
// 			// Depending on the time of day, this should be seen as today’s order.
// 			expectedCount: 2,
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
// 			name: "Valid statuses using TimeZoneOffset +0 (UTC) for store 1",
// 			filter: types.OrdersTimeZoneFilter{
// 				StoreID:        uintPtr(1),
// 				TimeZoneOffset: uintPtr(0),
// 			},
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

func TestOrderService_GeneratePDFReceipt(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	// Retrieve an order from store 1 to test.
	var order data.Order
	if err := db.Where("store_id = ?", 1).First(&order).Error; err != nil {
		t.Fatalf("Failed to fetch an order from store 1: %v", err)
	}

	pdfBytes, err := module.Service.GeneratePDFReceipt(order.ID)
	assert.NoError(t, err, "GeneratePDFReceipt should not return an error")
	assert.NotEmpty(t, pdfBytes, "Generated PDF should not be empty")
}

func TestOrderService_GenerateSuborderBarcodePDF(t *testing.T) {
	_ = resetTestData(t)
	module := tests.GetOrdersModule()

	suborderID := uint(1)
	// barcode text font might lead to a critical errors
	_ = utils.InitBarcodeFont()

	pdfBytes, err := module.Service.GenerateSuborderBarcodePDF(suborderID)
	assert.NoError(t, err, fmt.Sprintf("GenerateSuborderBarcodePDF returned an error for suborderID %d", suborderID))
	assert.NotEmpty(t, pdfBytes, "Generated PDF for suborder should not be empty")
}

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

	start := time.Now().Add(-72 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	filter := types.OrdersExportFilterQuery{
		StartDate: &start,
		EndDate:   &end,
		StoreID:   uintPtr(1),
	}
	exports, err := module.Service.ExportOrders(&filter)
	assert.NoError(t, err, "ExportOrders should not error")
	// Assuming test data for store 1 contains 2 orders.
	assert.Equal(t, 2, len(exports), "Expected 2 orders for export in store 1")
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

// func TestOrderService_CreateOrder_Combined(t *testing.T) {
// 	_ = resetTestData(t)
// 	module := tests.GetOrdersModule()

// 	err := censor.InitializeCensorForTests()
// 	if err != nil {
// 		t.Fatalf("Failed to initialize censor for tests: %v", err)
// 	}

// 	testCases := []struct {
// 		name              string
// 		dto               types.CreateOrderDTO
// 		storeID           uint
// 		expectedError     bool
// 		expectedErrSubstr string
// 		expectedTotal     float64
// 	}{
// 		{
// 			name: "Successful order creation - John Doe, quantity 2",
// 			dto: types.CreateOrderDTO{
// 				CustomerName: "John",
// 				Suborders: []types.CreateSubOrderDTO{
// 					{
// 						StoreProductSizeID: 1, // product size with store_price 2.75
// 						Quantity:           2,
// 						StoreAdditivesIDs:  []uint{1}, // additive with store_price 0.55
// 					},
// 				},
// 				StoreID: 1,
// 			},
// 			storeID:       1,
// 			expectedError: false,
// 			// Expected total = 2 * (2.75 + 0.55) = 6.60
// 			expectedTotal: 6.60,
// 		},
// 		{
// 			name: "Failure due to empty suborders",
// 			dto: types.CreateOrderDTO{
// 				CustomerName: "John",
// 				Suborders:    []types.CreateSubOrderDTO{},
// 				StoreID:      1,
// 			},
// 			storeID:           1,
// 			expectedError:     true,
// 			expectedErrSubstr: "order can not be empty",
// 		},
// 		{
// 			name: "Failure due to censored customer name",
// 			dto: types.CreateOrderDTO{
// 				CustomerName: "fuck", // Assume this is rejected by the censor validator.
// 				Suborders: []types.CreateSubOrderDTO{
// 					{
// 						StoreProductSizeID: 1,
// 						Quantity:           1,
// 						StoreAdditivesIDs:  []uint{1},
// 					},
// 				},
// 				StoreID: 1,
// 			},
// 			storeID:           1,
// 			expectedError:     true,
// 			expectedErrSubstr: "inappropriate", // Expect an error message indicating censorship.
// 		},
// 		{
// 			name: "Failure due to invalid product size",
// 			dto: types.CreateOrderDTO{
// 				CustomerName: "John",
// 				Suborders: []types.CreateSubOrderDTO{
// 					{
// 						StoreProductSizeID: 9999, // non-existent ID
// 						Quantity:           1,
// 						StoreAdditivesIDs:  []uint{1},
// 					},
// 				},
// 				StoreID: 1,
// 			},
// 			storeID:           1,
// 			expectedError:     true,
// 			expectedErrSubstr: "invalid store product size ID",
// 		},
// 	}

// 	// Run each test case.
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			order, err := module.Service.CreateOrder(tc.storeID, &tc.dto)
// 			if tc.expectedError {
// 				assert.Error(t, err, "Expected error in case %q", tc.name)
// 				if tc.expectedErrSubstr != "" {
// 					assert.Contains(t, err.Error(), tc.expectedErrSubstr, "Error message should contain %q", tc.expectedErrSubstr)
// 				}
// 				assert.Nil(t, order, "Order should be nil when creation fails")
// 			} else {
// 				assert.NoError(t, err, "Unexpected error in case %q", tc.name)
// 				assert.NotNil(t, order, "Returned order should not be nil in case %q", tc.name)
// 				// Verify total and status.
// 				assert.InDelta(t, tc.expectedTotal, order.Total, 0.01, "Order total is incorrect")
// 				assert.Equal(t, data.OrderStatusPending, order.Status, "New order should have PENDING status")
// 			}
// 		})
// 	}
// }

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

func TestOrderService_CompleteSubOrder_Combined(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	// Insert our custom test order.
	orderID, suborderIDs := insertTestOrderWithTwoSuborders(t, db)
	assert.Len(t, suborderIDs, 2, "Expected two suborders inserted")

	// Verify initial order status is PENDING.
	orderBefore := &data.Order{}
	err := db.Preload("Suborders").Where("id = ?", orderID).First(orderBefore).Error
	assert.NoError(t, err)
	assert.Equal(t, data.OrderStatusPending, orderBefore.Status, "Initial order status should be PENDING")

	// --- Subtest 1: Complete first suborder only.
	t.Run("Complete first suborder only", func(t *testing.T) {
		err := module.Service.CompleteSubOrder(orderID, suborderIDs[0])
		assert.NoError(t, err, "Completing first suborder should succeed")

		// Reload order status.
		orderAfterFirst := &data.Order{}
		err = db.Where("id = ?", orderID).First(orderAfterFirst).Error
		assert.NoError(t, err)
		// Not all suborders completed; status should remain PENDING.
		assert.Equal(t, data.OrderStatusPending, orderAfterFirst.Status, "Order status should remain PENDING")
	})

	// --- Subtest 2: Complete second suborder and check overall order update.
	t.Run("Complete all suborders", func(t *testing.T) {
		err := module.Service.CompleteSubOrder(orderID, suborderIDs[1])
		assert.NoError(t, err, "Completing second suborder should succeed")

		// Reload order status.
		orderAfterSecond := &data.Order{}
		err = db.Where("id = ?", orderID).First(orderAfterSecond).Error
		assert.NoError(t, err)
		// With both suborders complete, the order status should update.
		assert.Equal(t, data.OrderStatusCompleted, orderAfterSecond.Status, "Order status should update to COMPLETED")
	})

	// --- Subtest 3: Attempt to complete a non-existent suborder.
	t.Run("Complete non-existent suborder", func(t *testing.T) {
		err := module.Service.CompleteSubOrder(orderID, 99999) // Non-existent ID.
		assert.Error(t, err, "Should return error for non-existent suborder")
		assert.Contains(t, err.Error(), "failed to check suborder", "Error should mention suborder failure")
	})
}

func TestOrderService_CompleteSubOrderByBarcode(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	orderID, suborderIDs := insertTestOrderWithTwoSuborders(t, db)
	assert.Len(t, suborderIDs, 2, "Expected two suborders inserted")

	// Complete one suborder via barcode.
	dto, err := module.Service.CompleteSubOrderByBarcode(suborderIDs[0])
	assert.NoError(t, err, "Completing suborder by barcode should succeed")
	assert.NotNil(t, dto, "Returned suborder DTO should not be nil")

	// Order status remains PENDING because only one suborder is complete.
	order, err := module.Service.GetOrderById(orderID)
	assert.NoError(t, err)
	assert.Equal(t, data.OrderStatusPending, order.Status, "Order status should remain PENDING if not all suborders are complete")
}

func TestOrderService_AcceptSubOrder(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	// Insert an order with two suborders.
	_, suborderIDs := insertTestOrderWithTwoSuborders(t, db)
	assert.Len(t, suborderIDs, 2, "Expected two suborders inserted")

	// Accept the first suborder (transitions from PENDING to PREPARING).
	err := module.Service.AcceptSubOrder(suborderIDs[0])
	assert.NoError(t, err, "Accepting suborder should succeed for pending suborder")

	sub, err := module.Repo.GetSuborderByID(suborderIDs[0])
	assert.NoError(t, err)
	assert.Equal(t, data.SubOrderStatusPreparing, sub.Status, "Suborder status should be updated to PREPARING")

	// Trying to accept a non-pending suborder should error.
	err = module.Service.AcceptSubOrder(suborderIDs[0])
	assert.Error(t, err, "Accepting an already accepted suborder should error")
	assert.Contains(t, err.Error(), "is not pending", "Error message should indicate suborder is not pending")
}

func TestOrderService_AdvanceSubOrderStatus(t *testing.T) {
	db := resetTestData(t)
	module := tests.GetOrdersModule()

	// Insert an order with two suborders.
	orderID, suborderIDs := insertTestOrderWithTwoSuborders(t, db)
	assert.Len(t, suborderIDs, 2, "Expected two suborders inserted")

	// Advance the second suborder from PENDING to PREPARING.
	dto, err := module.Service.AdvanceSubOrderStatus(suborderIDs[1])
	assert.NoError(t, err, "Advancing suborder status should succeed")
	assert.Equal(t, data.SubOrderStatusPreparing, dto.Status, "Suborder should transition to PREPARING")

	// Advance the same suborder from PREPARING to COMPLETED.
	dto, err = module.Service.AdvanceSubOrderStatus(suborderIDs[1])
	assert.NoError(t, err, "Advancing suborder status again should succeed")
	assert.Equal(t, data.SubOrderStatusCompleted, dto.Status, "Suborder should transition to COMPLETED")

	// After both suborders are complete, the order status should update to COMPLETED.
	// For this test, complete the first suborder as well.
	err = module.Service.CompleteSubOrder(orderID, suborderIDs[0])
	assert.NoError(t, err, "Completing first suborder should succeed")
	order, err := module.Service.GetOrderById(orderID)
	assert.NoError(t, err)
	assert.Equal(t, data.OrderStatusCompleted, order.Status, "Order status should be COMPLETED when all suborders are complete")
}
