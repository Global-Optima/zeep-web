package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
)

type DeductedStoreInventory struct {
	StoreStocks     []data.StoreStock
	StoreProvisions []data.StoreProvision
}

type RecalculateInput struct {
	IngredientIDs   []uint
	ProvisionIDs    []uint
	ProductSizeIDs  []uint
	AdditiveIDs     []uint
	FrozenInventory *FrozenInventory // optional
}

type FrozenInventory struct {
	Ingredients map[uint]float64
	Provisions  map[uint]float64
}

func (f *FrozenInventory) GetIDs() *InventoryIDsList {
	inventoryIDsList := &InventoryIDsList{
		IngredientIDs: make([]uint, 0, len(f.Ingredients)),
		ProvisionIDs:  make([]uint, 0, len(f.Provisions)),
	}

	for id := range f.Ingredients {
		inventoryIDsList.IngredientIDs = append(inventoryIDsList.IngredientIDs, id)
	}

	for id := range f.Provisions {
		inventoryIDsList.ProvisionIDs = append(inventoryIDsList.ProvisionIDs, id)
	}

	return inventoryIDsList
}

type FrozenInventoryFilter struct {
	IngredientIDs []uint
	ProvisionIDs  []uint
}

type InventoryIDsList struct {
	IngredientIDs []uint
	ProvisionIDs  []uint
}

type InventoryUsage struct {
	Ingredients map[uint]float64
	Provisions  map[uint]float64
}

type DeductedInventoryMap struct {
	IngredientStoreStockMap     map[uint]*data.StoreStock
	ProvisionStoreProvisionsMap map[uint][]*data.StoreProvision
}

func (d *DeductedInventoryMap) GetIDs() *InventoryIDsList {
	inventoryIDsList := &InventoryIDsList{
		IngredientIDs: make([]uint, 0, len(d.IngredientStoreStockMap)),
		ProvisionIDs:  make([]uint, 0, len(d.ProvisionStoreProvisionsMap)),
	}

	for id := range d.IngredientStoreStockMap {
		inventoryIDsList.IngredientIDs = append(inventoryIDsList.IngredientIDs, id)
	}
	for id := range d.ProvisionStoreProvisionsMap {
		inventoryIDsList.ProvisionIDs = append(inventoryIDsList.ProvisionIDs, id)
	}

	return inventoryIDsList
}

type FrozenInventoryDTO struct {
	Ingredients []FrozenIngredientDTO `json:"ingredients"`
	Provisions  []FrozenProvisionDTO  `json:"provisions"`
}

type FrozenIngredientDTO struct {
	IngredientDTO ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity      float64                       `json:"quantity"`
}

type FrozenProvisionDTO struct {
	ProvisionDTO provisionsTypes.ProvisionDTO `json:"provision"`
	Volume       float64                      `json:"volume"`
}
