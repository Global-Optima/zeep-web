package sku

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku/types"
)

type SKUService interface {
	GetAllSKUs(filter *types.SKUFilter) ([]types.SKUResponse, error)
	GetSKUByID(skuID uint) (*types.SKUResponse, error)
	CreateSKU(req *types.CreateSKURequest) (*types.SKUResponse, error)
	UpdateSKU(skuID uint, req *types.UpdateSKURequest) (*types.SKUResponse, error)
	DeleteSKU(skuID uint) error
	DeactivateSKU(skuID uint) error
}

type skuService struct {
	repo SKURepository
}

func NewSKUService(repo SKURepository) SKUService {
	return &skuService{repo: repo}
}

func (s *skuService) GetAllSKUs(filter *types.SKUFilter) ([]types.SKUResponse, error) {
	skus, err := s.repo.GetAllSKUs(filter)
	if err != nil {
		return nil, err
	}

	var skuResponses []types.SKUResponse
	for _, sku := range skus {
		skuResponses = append(skuResponses, *types.ConvertSKUToSKUResponse(&sku))
	}

	return skuResponses, nil
}

func (s *skuService) GetSKUByID(skuID uint) (*types.SKUResponse, error) {
	sku, err := s.repo.GetSKUByID(skuID)
	if err != nil {
		return nil, err
	}

	if sku == nil {
		return nil, errors.New("SKU not found")
	}

	skuResponse := types.ConvertSKUToSKUResponse(sku)

	return skuResponse, nil
}

func (s *skuService) CreateSKU(req *types.CreateSKURequest) (*types.SKUResponse, error) {
	sku := types.ConvertCreateSKURequestToSKU(req)

	err := s.repo.CreateSKU(sku)
	if err != nil {
		return nil, err
	}

	skuResponse := types.ConvertSKUToSKUResponse(sku)
	return skuResponse, nil
}

func (s *skuService) UpdateSKU(skuID uint, req *types.UpdateSKURequest) (*types.SKUResponse, error) {
	updateFields := make(map[string]interface{})

	err := types.ConvertUpdateSKURequestToMap(req, updateFields)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateSKUFields(skuID, updateFields)
	if err != nil {
		return nil, err
	}
	updatedSKU := types.ConvertSKUToSKUResponse(updated)
	return updatedSKU, nil
}

func (s *skuService) DeleteSKU(skuID uint) error {
	return s.repo.DeleteSKU(skuID)
}

func (s *skuService) DeactivateSKU(skuID uint) error {
	return s.repo.DeactivateSKU(skuID)
}
