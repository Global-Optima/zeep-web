package storeInventoryManagers

type recalculateContext struct {
	hasIngredients  bool
	hasProvisions   bool
	hasProductSizes bool
	hasAdditives    bool

	productSizesIngredientIDs []uint
	productSizesProvisionIDs  []uint
	additiveIngredientIDs     []uint
	additiveProvisionIDs      []uint

	storeProductIDsFromPS          []uint
	storeProductIDsFromIngredients []uint
	storeProductIDsFromProvisions  []uint

	storeAdditiveIDsFromAdditives   []uint
	storeAdditiveIDsFromIngredients []uint
	storeAdditiveIDsFromProvisions  []uint

	totalIngredientIDs []uint
	totalProvisionIDs  []uint
}

type additiveIngredientUsageRow struct {
	StoreAdditiveID  uint
	IngredientID     uint
	RequiredQuantity float64
}

type additiveProvisionUsageRow struct {
	StoreAdditiveID uint
	ProvisionID     uint
	RequiredVolume  float64
}

type productSizeIngredientRow struct {
	StoreProductSizeID uint
	IngredientID       uint
	RequiredQuantity   float64
}

type productSizeProvisionRow struct {
	StoreProductSizeID uint
	ProvisionID        uint
	RequiredVolume     float64
}

type usageKey struct {
	StoreProductSizeID uint
	ResourceID         uint // ingredientID or provisionID
}

// Aggregated usage for each storeProductSize -> resource.
type aggregatedUsage struct {
	Ingredient map[usageKey]float64
	Provision  map[usageKey]float64
}
