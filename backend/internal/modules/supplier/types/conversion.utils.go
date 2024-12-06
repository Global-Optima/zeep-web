package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToSupplierResponse(supplier data.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:           supplier.ID,
		Name:         supplier.Name,
		ContactEmail: supplier.ContactEmail,
		ContactPhone: supplier.ContactPhone,
		Address:      supplier.Address,
		CreatedAt:    supplier.CreatedAt,
		UpdatedAt:    supplier.UpdatedAt,
	}
}

func ToSupplier(dto CreateSupplierDTO) data.Supplier {
	return data.Supplier{
		Name:         dto.Name,
		ContactEmail: dto.ContactEmail,
		ContactPhone: dto.ContactPhone,
		Address:      dto.Address,
	}
}

func ConvertUpdateSupplierDTOToMap(dto UpdateSupplierDTO) (map[string]interface{}, error) {
	updateFields := make(map[string]interface{})

	if dto.Name != nil {
		updateFields["name"] = *dto.Name
	}
	if dto.ContactEmail != nil {
		updateFields["contact_email"] = *dto.ContactEmail
	}
	if dto.ContactPhone != nil {
		updateFields["contact_phone"] = *dto.ContactPhone
	}
	if dto.Address != nil {
		updateFields["address"] = *dto.Address
	}

	return updateFields, nil
}
