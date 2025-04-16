package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

type DeductedStoreInventory struct {
	StoreStocks     []data.StoreStock
	StoreProvisions []data.StoreProvision
}

type RecalculateInput struct {
	IngredientIDs   []uint
	ProvisionIDs    []uint
	ProductSizeIDs  []uint
	AdditiveIDs     []uint
	FrozenInventory *FrozenInventory //optional
}

type FrozenInventory struct {
	Ingredients map[uint]float64
	Provisions  map[uint]float64
}

type FrozenInventoryFilter struct {
	IngredientIDs []uint
	ProvisionIDs  []uint
}

type InventoryIDsLists struct {
	IngredientIDs []uint
	ProvisionIDs  []uint
}
