package supplier

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
)

type SupplierService interface {
	CreateSupplier(dto types.CreateSupplierDTO) (types.SupplierResponse, error)
	GetSupplierByID(id uint) (types.SupplierResponse, error)
	UpdateSupplier(id uint, dto types.UpdateSupplierDTO) error
	DeleteSupplier(id uint) error
	GetSuppliers() ([]types.SupplierResponse, error)
}

type supplierService struct {
	repo SupplierRepository
}

func NewSupplierService(repo SupplierRepository) SupplierService {
	return &supplierService{repo}
}

func (s *supplierService) CreateSupplier(dto types.CreateSupplierDTO) (types.SupplierResponse, error) {
	supplier := types.ToSupplier(dto)
	if err := s.repo.CreateSupplier(&supplier); err != nil {
		return types.SupplierResponse{}, err
	}
	return types.ToSupplierResponse(supplier), nil
}

func (s *supplierService) GetSupplierByID(id uint) (types.SupplierResponse, error) {
	supplier, err := s.repo.GetSupplierByID(id)
	if err != nil {
		return types.SupplierResponse{}, err
	}
	return types.ToSupplierResponse(*supplier), nil
}

func (s *supplierService) UpdateSupplier(id uint, dto types.UpdateSupplierDTO) error {
	supplier, err := s.repo.GetSupplierByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch supplier with ID %d: %w", id, err)
	}

	if dto.Name != nil {
		supplier.Name = *dto.Name
	}
	if dto.ContactEmail != nil {
		supplier.ContactEmail = *dto.ContactEmail
	}
	if dto.ContactPhone != nil {
		supplier.ContactPhone = *dto.ContactPhone
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

func (s *supplierService) GetSuppliers() ([]types.SupplierResponse, error) {
	suppliers, err := s.repo.ListSuppliers()
	if err != nil {
		return nil, err
	}
	responses := make([]types.SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = types.ToSupplierResponse(supplier)
	}
	return responses, nil
}
