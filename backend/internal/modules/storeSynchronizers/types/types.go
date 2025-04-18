package types

type UnsyncData struct {
	ProductSizeIDs []uint
	AdditiveIDs    []uint
	IngredientIDs  []uint
	ProvisionIDs   []uint
}

type SynchronizationStatus struct {
	IsSync       bool   `json:"isSync"`
	LastSyncDate string `json:"lastSyncDate"`
}
