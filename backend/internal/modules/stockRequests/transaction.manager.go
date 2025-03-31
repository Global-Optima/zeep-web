package stockRequests

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"gorm.io/gorm"
)

type TransactionManager interface {
	HandleCompleteStockRequest(request *data.StockRequest) error
}

type transactionManager struct {
	db                *gorm.DB
	repo              StockRequestRepository
	stockMaterialRepo stockMaterial.StockMaterialRepository
}

func NewTransactionManager(
	db *gorm.DB,
	repo StockRequestRepository,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
) TransactionManager {
	return &transactionManager{
		db:                db,
		repo:              repo,
		stockMaterialRepo: stockMaterialRepo,
	}
}

func (m *transactionManager) HandleCompleteStockRequest(request *data.StockRequest) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		stockMaterialIDs := make([]uint, len(request.Ingredients))
		repoTx := m.repo.CloneWithTransaction(tx)

		for i, ingredient := range request.Ingredients {
			stockMaterialIDs[i] = ingredient.StockMaterialID

			dates := types.UpdateIngredientDates{
				DeliveredDate:  time.Now(),
				ExpirationDate: time.Now().AddDate(0, 0, ingredient.StockMaterial.ExpirationPeriodInDays),
			}

			if err := repoTx.UpdateStockRequestIngredientDates(ingredient.ID, &dates); err != nil {
				return fmt.Errorf("failed to update ingredient dates for stock material ID %d: %w", ingredient.StockMaterialID, err)
			}

			if err := repoTx.UpsertToStoreStock(request.StoreID, ingredient.StockMaterialID, ingredient.Quantity); err != nil {
				return fmt.Errorf("failed to update store warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
			}
		}

		stockMaterials, err := m.stockMaterialRepo.GetStockMaterialsByIDs(stockMaterialIDs)
		if err != nil {
			return fmt.Errorf("failed to fetch stock materials: %w", err)
		}

		ingredientIDs := make([]uint, len(stockMaterials))
		for i, sm := range stockMaterials {
			ingredientIDs[i] = sm.IngredientID
		}

		request.Status = data.StockRequestCompleted
		if err := repoTx.UpdateStockRequestStatus(request); err != nil {
			return fmt.Errorf("failed to update stock request status: %w", err)
		}

		if err := data.RecalculateOutOfStock(tx, request.StoreID, ingredientIDs, nil, nil); err != nil {
			return err
		}

		return nil
	})
}
