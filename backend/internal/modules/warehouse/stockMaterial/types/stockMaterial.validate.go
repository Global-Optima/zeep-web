package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateAndApplyUpdate(stockMaterial *data.StockMaterial, req *UpdateStockMaterialDTO) (*data.StockMaterial, error) {
	if stockMaterial == nil {
		return nil, errors.New("stockMaterial cannot be nil")
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("stockMaterial name cannot be empty")
		}
		stockMaterial.Name = *req.Name
	}

	if req.Description != nil {
		stockMaterial.Description = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return nil, errors.New("stockMaterial safety stock must be greater than zero")
		}
		stockMaterial.SafetyStock = *req.SafetyStock
	}

	if req.UnitID != nil {
		stockMaterial.UnitID = *req.UnitID
	}

	if req.CategoryID != nil {
		stockMaterial.CategoryID = *req.CategoryID
	}

	if req.IngredientID != nil {
		stockMaterial.IngredientID = *req.IngredientID
	}

	if req.Barcode != nil {
		stockMaterial.Barcode = *req.Barcode
	}

	if req.ExpirationPeriodInDays != nil {
		stockMaterial.ExpirationPeriodInDays = *req.ExpirationPeriodInDays
	}

	if req.IsActive != nil {
		stockMaterial.IsActive = *req.IsActive
	}

	stockMaterial.UpdatedAt = time.Now()

	return stockMaterial, nil
}

func ValidatePackageUpdates(pkg *data.StockMaterialPackage, dto *UpdateStockMaterialPackagesDTO) error {
	if dto.Size != nil && *dto.Size <= 0 {
		return fmt.Errorf("invalid size for package ID %d", *dto.StockMaterialPackageID)
	}
	if dto.UnitID != nil && *dto.UnitID == 0 {
		return fmt.Errorf("invalid unit ID for package ID %d", *dto.StockMaterialPackageID)
	}
	return nil
}

func ValidatePackageDTO(pkgDTO UpdateStockMaterialPackagesDTO) error {
	if pkgDTO.StockMaterialPackageID == nil && pkgDTO.Size == nil && pkgDTO.UnitID == nil {
		return fmt.Errorf("empty package object is not allowed")
	}

	if pkgDTO.StockMaterialPackageID == nil {
		if pkgDTO.Size == nil || pkgDTO.UnitID == nil {
			return fmt.Errorf("new packages must have both size and unitId")
		}
	}

	return nil
}

func ValidatePackageUpdates(pkg *data.StockMaterialPackage, dto *UpdateStockMaterialPackagesDTO) error {
	if dto.Size != nil && *dto.Size <= 0 {
		return fmt.Errorf("invalid size for package ID %d", *dto.StockMaterialPackageID)
	}
	if dto.UnitID != nil && *dto.UnitID == 0 {
		return fmt.Errorf("invalid unit ID for package ID %d", *dto.StockMaterialPackageID)
	}
	return nil
}
