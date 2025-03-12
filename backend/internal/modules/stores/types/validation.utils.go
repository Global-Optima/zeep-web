package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	facilityAddressesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/facilityAddresses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"
)

type StoreUpdateModels struct {
	Store           *data.Store
	FacilityAddress *data.FacilityAddress
}

func UpdateStoreFields(dto *UpdateStoreDTO, store *data.Store, facilityAddress *data.FacilityAddress) (*StoreUpdateModels, error) {
	store.Warehouse = data.Warehouse{}
	store.Franchisee = nil

	if dto.Name != "" {
		store.Name = dto.Name
	}

	store.FranchiseeID = dto.FranchiseeID

	if dto.WarehouseID != nil {
		store.WarehouseID = *dto.WarehouseID
	}
	if dto.IsActive != nil {
		store.IsActive = *dto.IsActive
	}
	if dto.ContactPhone != "" {
		if !utils.IsValidPhone(dto.ContactPhone, utils.DEFAULT_PHONE_NUMBER_REGION) {
			return nil, moduleErrors.ErrValidation.WithDetails("phoneNumber")
		}
		store.ContactPhone = dto.ContactPhone
	}
	if dto.ContactEmail != "" {
		if !utils.IsValidEmail(dto.ContactEmail) {
			return nil, moduleErrors.ErrValidation.WithDetails("email")
		}
		store.ContactEmail = dto.ContactEmail
	}
	if dto.StoreHours != "" {
		store.StoreHours = dto.StoreHours
	}
	if dto.FacilityAddress != nil && dto.FacilityAddress.Address != "" {
		facilityAddress = facilityAddressesTypes.MapToFacilityAddressModel(dto.FacilityAddress, facilityAddress)
	}

	logrus.Info(*store)
	return &StoreUpdateModels{
		Store:           store,
		FacilityAddress: facilityAddress,
	}, nil
}

func CreateStoreFields(dto *CreateStoreDTO) (*data.Store, error) {
	store := &data.Store{}

	store.Name = dto.Name
	store.FranchiseeID = dto.FranchiseeID
	store.WarehouseID = dto.WarehouseID
	store.FacilityAddress = *facilityAddressesTypes.MapToFacilityAddressModel(&dto.FacilityAddress, nil)

	if !utils.IsValidPhone(dto.ContactPhone, utils.DEFAULT_PHONE_NUMBER_REGION) {
		return nil, moduleErrors.ErrValidation.WithDetails("phoneNumber")
	}
	store.ContactPhone = dto.ContactPhone

	if !utils.IsValidEmail(dto.ContactEmail) {
		return nil, moduleErrors.ErrValidation.WithDetails("email")
	}
	store.ContactEmail = dto.ContactEmail

	store.StoreHours = dto.StoreHours

	return store, nil
}
