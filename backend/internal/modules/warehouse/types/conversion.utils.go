package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func ConvertToListStoresResponse(stores []data.Store) []ListStoresResponse {
	response := make([]ListStoresResponse, len(stores))
	for i, store := range stores {
		response[i] = ListStoresResponse{
			StoreID: store.ID,
			Name:    store.Name,
			Status:  store.Status,
		}
	}
	return response
}
