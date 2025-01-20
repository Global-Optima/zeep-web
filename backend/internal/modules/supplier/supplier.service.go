package supplier

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
)

type SupplierService interface {
	CreateSupplier(dto types.CreateSupplierDTO) (*types.SupplierResponse, error)
	GetSupplierByID(id uint) (types.SupplierResponse, error)
	UpdateSupplier(id uint, dto types.UpdateSupplierDTO) error
	DeleteSupplier(id uint) error
	GetSuppliers(filter types.SuppliersFilter) ([]types.SupplierResponse, error)

	UpsertMaterialsForSupplier(supplierID uint, dto types.UpsertSupplierMaterialsDTO) error
	GetMaterialsBySupplier(supplierID uint) ([]types.SupplierMaterialResponse, error)
}

type supplierService struct {
	repo SupplierRepository
}

func NewSupplierService(repo SupplierRepository) SupplierService {
	return &supplierService{repo}
}

func (s *supplierService) CreateSupplier(dto types.CreateSupplierDTO) (*types.SupplierResponse, error) {
	exists, err := s.repo.ExistsByContactPhone(dto.ContactPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to check supplier existence: %w", err)
	}
	if exists {
		return nil, errors.New("a supplier with this contact phone already exists")
	}

	supplier := types.ToSupplier(dto)
	if err := s.repo.CreateSupplier(&supplier); err != nil {
		return nil, err
	}

	response := types.ToSupplierResponse(supplier)
	return &response, nil
}

func (s *supplierService) GetSupplierByID(id uint) (types.SupplierResponse, error) {
	supplier, err := s.repo.GetSupplierByID(id)
	if err != nil {
		return types.SupplierResponse{}, err
	}
	return types.ToSupplierResponse(*supplier), nil
}

func (s *supplierService) UpdateSupplier(id uint, dto types.UpdateSupplierDTO) error {
	supplier := &data.Supplier{}
	if dto.Name != nil {
		supplier.Name = *dto.Name
	}
	if dto.ContactEmail != nil {
		supplier.ContactEmail = *dto.ContactEmail
	}
	if dto.ContactPhone != nil {
		supplier.ContactPhone = *dto.ContactPhone
	}
	if dto.City != nil {
		supplier.City = *dto.City
	}
	if dto.Address != nil {
		supplier.Address = *dto.Address
	}

	if err := s.repo.UpdateSupplier(id, supplier); err != nil {
		return fmt.Errorf("failed to update supplier with ID %d: %w", id, err)
	}

	return nil
}

func (s *supplierService) DeleteSupplier(id uint) error {
	return s.repo.DeleteSupplier(id)
}

func (s *supplierService) GetSuppliers(filter types.SuppliersFilter) ([]types.SupplierResponse, error) {
	suppliers, err := s.repo.GetAllSuppliers(filter)
	if err != nil {
		return nil, err
	}
	responses := make([]types.SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = types.ToSupplierResponse(supplier)
	}
	return responses, nil
}

func (s *supplierService) UpsertMaterialsForSupplier(supplierID uint, dto types.UpsertSupplierMaterialsDTO) error {
	materials := make([]data.SupplierMaterial, len(dto.Materials))
	for i, materialDTO := range dto.Materials {
		materials[i] = data.SupplierMaterial{
			SupplierID:      supplierID,
			StockMaterialID: materialDTO.StockMaterialID,
			SupplierPrices: []data.SupplierPrice{
				{
					BasePrice:     materialDTO.BasePrice,
					EffectiveDate: time.Now(),
				},
			},
		}
	}

	if err := s.repo.UpsertSupplierMaterials(supplierID, materials); err != nil {
		return fmt.Errorf("failed to upsert materials: %w", err)
	}

	return nil
}

func (s *supplierService) GetMaterialsBySupplier(supplierID uint) ([]types.SupplierMaterialResponse, error) {
	materials, err := s.repo.GetMaterialsBySupplier(supplierID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch materials: %w", err)
	}

	return types.ToSupplierMaterialResponses(materials), nil
}
