package orders

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TransactionManager interface {
	SetNextSubOrderStatus(suborder *data.Suborder) error
}

type transactionManager struct {
	db                        *gorm.DB
	repo                      OrderRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
	notificationService       notifications.NotificationService
	logger                    *zap.SugaredLogger
}

func NewTransactionManager(
	db *gorm.DB,
	repo OrderRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) TransactionManager {
	return &transactionManager{
		db:                        db,
		repo:                      repo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
		notificationService:       notificationService,
		logger:                    logger,
	}
}

func (m *transactionManager) SetNextSubOrderStatus(suborder *data.Suborder) error {
	if suborder == nil {
		return fmt.Errorf("suborder ID is nil")
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.repo.CloneWithTransaction(tx)
		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		// Attempt to advance suborder status
		if err := m.nextSuborderStatus(&repoTx, storeInventoryManagerRepoTx, suborder.ID, suborder); err != nil {
			return err
		}

		// Sync and update order status
		if err := m.updateOrderStatusBySuborder(&repoTx, suborder.ID); err != nil {
			return err
		}

		return nil
	}) // Handle fallback if suborder is already completed within time gap
}

func (m *transactionManager) nextSuborderStatus(repoTx OrderRepository, storeInventoryManagerRepoTx storeInventoryManagers.StoreInventoryManagerRepository, subOrderID uint, suborder *data.Suborder) error {
	currentStatus := suborder.Status
	nextStatus, ok := allowedTransitions[currentStatus]
	if !ok {
		return fmt.Errorf("no allowed transition from status %s", currentStatus)
	}

	completedAt := time.Now()
	update := types.UpdateSubOrderDTO{
		Status:      nextStatus,
		CompletedAt: &completedAt,
	}
	if err := repoTx.UpdateSubOrderStatus(subOrderID, update); err != nil {
		return fmt.Errorf("failed to update suborder status: %w", err)
	}

	// If suborder is completed, deduct ingredients
	if nextStatus == data.SubOrderStatusCompleted {
		if err := m.handleSuborderCompletion(repoTx, storeInventoryManagerRepoTx, subOrderID); err != nil {
			return err
		}
	}

	return nil
}

func (m *transactionManager) handleSuborderCompletion(repoTx OrderRepository, storeInventoryManagerRepoTx storeInventoryManagers.StoreInventoryManagerRepository, subOrderID uint) error {
	suborder, err := repoTx.GetSuborderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve updated suborder: %w", err)
	}

	order, err := repoTx.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	inventoryMap := &types.DeductedInventoryMap{
		IngredientStoreStockMap:     make(map[uint]*data.StoreStock),
		ProvisionStoreProvisionsMap: make(map[uint][]data.StoreProvision),
	}

	if err := m.deductSuborderInventoryFromStock(storeInventoryManagerRepoTx, order.StoreID, suborder, inventoryMap); err != nil {
		return fmt.Errorf("failed to deduct ingredients: %w", err)
	}

	//filter storeStocks for recalculation: only keep those below or equal to lowStockThreshold
	for id, stock := range inventoryMap.IngredientStoreStockMap {
		if stock.Quantity > stock.LowStockThreshold {
			delete(inventoryMap.IngredientStoreStockMap, id)
		}
	}

	ingredientIDs, provisionIDs := inventoryMap.GetKeys()
	err = storeInventoryManagerRepoTx.RecalculateStoreInventory(order.StoreID, &storeInventoryManagersTypes.RecalculateInput{
		IngredientIDs: ingredientIDs,
		ProvisionIDs:  provisionIDs,
	})

	if err != nil {
		return fmt.Errorf("failed to recalculate out of stock: %w", err)
	}

	m.notifyLowStockIngredients(order, inventoryMap.IngredientStoreStockMap)
	return nil
}

func (m *transactionManager) updateOrderStatusBySuborder(repoTx OrderRepository, subOrderID uint) error {
	order, err := repoTx.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := repoTx.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	hasPreparing, allCompleted := m.evaluateSuborderStatuses(suborders)

	switch {
	case hasPreparing:
		return m.ensureOrderStatus(repoTx, order, data.OrderStatusPreparing, nil)

	case allCompleted:
		newStatus := data.OrderStatusCompleted
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		}
		now := time.Now()
		return m.ensureOrderStatus(repoTx, order, newStatus, &now)

	default:
		return m.ensureOrderStatus(repoTx, order, data.OrderStatusPreparing, nil)
	}
}

func (m *transactionManager) ensureOrderStatus(repoTx OrderRepository, order *data.Order, status data.OrderStatus, completedAt *time.Time) error {
	if order.Status == status {
		return nil
	}

	update := types.UpdateOrderDTO{
		Status:      status,
		CompletedAt: completedAt,
	}
	if err := repoTx.UpdateOrderStatus(order.ID, update); err != nil {
		return fmt.Errorf("failed to update order status to %s: %w", status, err)
	}
	return nil
}

func (m *transactionManager) evaluateSuborderStatuses(suborders []data.Suborder) (hasPreparing bool, allCompleted bool) {
	hasPreparing = false
	allCompleted = true
	for _, so := range suborders {
		if so.Status == data.SubOrderStatusPreparing {
			hasPreparing = true
		}
		if so.Status != data.SubOrderStatusCompleted {
			allCompleted = false
		}
	}
	return
}

func (m *transactionManager) deductSuborderInventoryFromStock(storeInventoryManagerRepoTx storeInventoryManagers.StoreInventoryManagerRepository, storeID uint, suborder *data.Suborder, inventoryMap *types.DeductedInventoryMap) error {
	if storeID == 0 || suborder == nil || inventoryMap == nil {
		return fmt.Errorf("failed to deduct suborder: invalid input parameters passed")
	}

	deductedInventoryFromProductSize, err := storeInventoryManagerRepoTx.DeductStoreInventoryByProductSize(storeID, suborder.StoreProductSize.ProductSizeID)
	if err != nil {
		return fmt.Errorf("failed to deduct product size ingredients: %w", err)
	}
	for _, stock := range deductedInventoryFromProductSize.StoreStocks {
		if existingStock, exists := inventoryMap.IngredientStoreStockMap[stock.IngredientID]; exists {
			existingStock.Quantity = stock.Quantity
		} else {
			inventoryMap.IngredientStoreStockMap[stock.IngredientID] = &stock
		}
	}

	for _, subAdditive := range suborder.SuborderAdditives {
		deductedInventoryFromAdditive, err := storeInventoryManagerRepoTx.DeductStoreInventoryByAdditive(storeID, subAdditive.StoreAdditive.AdditiveID)
		if err != nil {
			return fmt.Errorf("failed to deduct additive ingredients: %w", err)
		}
		for _, stock := range deductedInventoryFromAdditive.StoreStocks {
			if existingStock, exists := inventoryMap.IngredientStoreStockMap[stock.IngredientID]; exists {
				existingStock.Quantity = stock.Quantity
			} else {
				inventoryMap.IngredientStoreStockMap[stock.IngredientID] = &stock
			}
		}
	}

	return nil
}

func (m *transactionManager) notifyLowStockIngredients(order *data.Order, stockMap map[uint]*data.StoreStock) {
	for _, stock := range stockMap {
		if stock.Quantity <= stock.LowStockThreshold {
			notificationDetails := &details.StoreWarehouseRunOutDetails{
				BaseNotificationDetails: details.BaseNotificationDetails{
					ID:           order.StoreID,
					FacilityName: order.Store.Name,
				},
				StockItem:   stock.Ingredient.Name,
				StockItemID: stock.IngredientID,
			}
			err := m.notificationService.NotifyStoreWarehouseRunOut(notificationDetails)
			if err != nil {
				m.logger.Errorf("failed to send notification: %v", err)
			}
		}
	}
}
