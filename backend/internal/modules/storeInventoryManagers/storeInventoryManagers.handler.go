package storeInventoryManagers

type StoreInventoryManagerHandler struct {
	service StoreInventoryManagerService
}

func NewStoreInventoryManagerHandler(service StoreInventoryManagerService) *StoreInventoryManagerHandler {
	return &StoreInventoryManagerHandler{
		service: service,
	}
}
