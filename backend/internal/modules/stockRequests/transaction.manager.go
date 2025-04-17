package stockRequests

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"gorm.io/gorm"
)

type TransactionManager interface {
	HandleCompleteStockRequest(request *data.StockRequest) error
	HandleAcceptedWithChange(
		request *data.StockRequest,
		storeID uint,
		items []types.StockRequestStockMaterialDTO,
		comment *string,
	) error
}

type transactionManager struct {
	db                        *gorm.DB
	repo                      StockRequestRepository
	stockMaterialRepo         stockMaterial.StockMaterialRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
}

func NewTransactionManager(
	db *gorm.DB,
	repo StockRequestRepository,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
) TransactionManager {
	return &transactionManager{
		db:                        db,
		repo:                      repo,
		stockMaterialRepo:         stockMaterialRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
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

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		err = storeInventoryManagerRepoTx.RecalculateStoreInventory(
			request.StoreID,
			&storeInventoryManagersTypes.RecalculateInput{
				IngredientIDs: ingredientIDs,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *transactionManager) HandleAcceptedWithChange(request *data.StockRequest, storeID uint, items []types.StockRequestStockMaterialDTO, comment *string) error {
	updatedIngredients := []data.StockRequestIngredient{}
	var changeDetails []types.StockRequestDetails

	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.repo.CloneWithTransaction(tx)
		stockMaterialIDs := make([]uint, len(items))
		for i, item := range items {
			stockMaterialIDs[i] = item.StockMaterialID
		}

		stockMaterials, err := m.stockMaterialRepo.GetStockMaterialsByIDs(stockMaterialIDs)
		if err != nil {
			return fmt.Errorf("failed to fetch stock materials: %w", err)
		}

		materialMap := make(map[uint]data.StockMaterial)
		ingredientIDs := make([]uint, len(stockMaterials))
		for i, sm := range stockMaterials {
			ingredientIDs[i] = sm.IngredientID
			materialMap[sm.ID] = sm
		}

		for _, item := range items {
			originalIngredient := findOriginalIngredient(request.Ingredients, item.StockMaterialID)

			if originalIngredient != nil {
				if originalIngredient.Quantity != item.Quantity {
					requestDetails := types.StockRequestDetails{
						OriginalMaterialName: originalIngredient.StockMaterial.Name,
						Quantity:             originalIngredient.Quantity,
						ActualQuantity:       item.Quantity,
					}
					changeDetails = append(changeDetails, requestDetails)

					// If accepted quantity is lower, return the difference to the warehouse.
					if originalIngredient.Quantity > item.Quantity {
						diff := originalIngredient.Quantity - item.Quantity
						_, err := repoTx.ReturnWarehouseStock(item.StockMaterialID, request.WarehouseID, diff)
						if err != nil {
							return fmt.Errorf("failed to return excess stock for material ID %d: %w", item.StockMaterialID, err)
						}
					}
				}
			} else {
				requestDetails := types.StockRequestDetails{
					MaterialName:   materialMap[item.StockMaterialID].Name,
					ActualQuantity: item.Quantity,
				}
				changeDetails = append(changeDetails, requestDetails)
			}

			if item.Quantity > 0 {
				if err := repoTx.UpsertToStoreStock(storeID, item.StockMaterialID, item.Quantity); err != nil {
					return fmt.Errorf("failed to add stock to store warehouse for stock material ID %d: %w", item.StockMaterialID, err)
				}
			}

			updatedIngredients = append(updatedIngredients, data.StockRequestIngredient{
				StockRequestID:  request.ID,
				StockMaterialID: item.StockMaterialID,
				Quantity:        item.Quantity,
				DeliveredDate:   time.Now(),
				ExpirationDate:  time.Now().AddDate(0, 0, materialMap[item.StockMaterialID].ExpirationPeriodInDays),
			})
		}

		if len(changeDetails) > 0 {
			err := repoTx.AddDetails(request.ID, changeDetails)
			if err != nil {
				return fmt.Errorf("failed to add details of changes for request ID %d: %w", request.ID, err)
			}
		}

		if err := repoTx.ReplaceStockRequestIngredients(*request, updatedIngredients); err != nil {
			return fmt.Errorf("failed to replace ingredients for stock request ID %d: %w", request.ID, err)
		}

		if comment != nil {
			if err := repoTx.AddStoreComment(request.ID, *comment); err != nil {
				return fmt.Errorf("failed to add store comment for request ID %d: %w", request.ID, err)
			}
		}

		request.Status = data.StockRequestAcceptedWithChange
		if err := repoTx.UpdateStockRequestStatus(request); err != nil {
			return fmt.Errorf("failed to update stock request status: %w", err)
		}

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		_ = storeInventoryManagerRepoTx.RecalculateStoreInventory(
			request.StoreID,
			&storeInventoryManagersTypes.RecalculateInput{
				IngredientIDs: ingredientIDs,
			},
		)

		return nil
	})
}
