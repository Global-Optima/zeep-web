package stores

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
)

type StoreService interface {
	GetAllStores() ([]types.StoreDTO, error)
	GetStoreEmployees(storeID uint) ([]types.EmployeeDTO, error)
}

type storeService struct {
	repo StoreRepository
}

func NewStoreService(repo StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) GetAllStores() ([]types.StoreDTO, error) {
	stores, err := s.repo.GetAllStores()
	if err != nil {
		return nil, err
	}

	storesDTOs := make([]types.StoreDTO, len(stores))
	for i, store := range stores {
		storesDTOs[i] = mapToStoreDTO(store)
	}

	return storesDTOs, nil
}

func (s *storeService) GetStoreEmployees(storeID uint) ([]types.EmployeeDTO, error) {
	employees, err := s.repo.GetStoreEmployees(storeID)
	if err != nil {
		return nil, err
	}

	employeeDTOs := make([]types.EmployeeDTO, len(employees))
	for i, employee := range employees {
		employeeDTOs[i] = mapToEmployeeDTO(employee)
	}

	return employeeDTOs, nil
}

func mapToStoreDTO(store data.Store) types.StoreDTO {
	return types.StoreDTO{
		ID:          store.ID,
		Name:        store.Name,
		IsFranchise: store.IsFranchise,
		FacilityAddress: &types.FacilityAddressDTO{
			ID:      store.FacilityAddress.ID,
			Address: store.FacilityAddress.Address,
		},
	}
}

func mapToEmployeeDTO(employee data.Employee) types.EmployeeDTO {
	var roleDTO *types.EmployeeRoleDTO

	if employee.Role != nil {
		roleDTO = &types.EmployeeRoleDTO{
			ID:   employee.Role.ID,
			Name: employee.Role.Name,
		}
	}

	return types.EmployeeDTO{
		ID:       employee.ID,
		Name:     employee.Name,
		Phone:    employee.Phone,
		Email:    employee.Email,
		IsActive: employee.IsActive,
		Role:     roleDTO,
	}
}
