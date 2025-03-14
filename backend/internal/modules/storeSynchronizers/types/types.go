package types

type UnsyncData struct {
	AdditiveIDs   []uint
	IngredientIDs []uint
}

type SynchronizationStatus struct {
	IsSync       bool   `json:"isSync"`
	LastSyncDate string `json:"lastSyncDate"`
}
