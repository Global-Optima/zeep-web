package types

import "errors"

func ValidateCreateSupplierDTO(input CreateSupplierDTO) error {
	if input.Name == "" {
		return errors.New("name is required")
	}

	if len(input.Name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}

	if input.ContactEmail != "" && len(input.ContactEmail) > 255 {
		return errors.New("contact_email must not exceed 255 characters")
	}

	if input.ContactPhone != "" && len(input.ContactPhone) > 20 {
		return errors.New("contact_phone must not exceed 20 characters")
	}

	if input.Address != "" && len(input.Address) > 255 {
		return errors.New("address must not exceed 255 characters")
	}

	return nil
}

func ValidateUpdateSupplierDTO(input UpdateSupplierDTO) error {
	if input.Name != nil && len(*input.Name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}

	if input.ContactEmail != nil && len(*input.ContactEmail) > 255 {
		return errors.New("contact_email must not exceed 255 characters")
	}

	if input.ContactPhone != nil && len(*input.ContactPhone) > 20 {
		return errors.New("contact_phone must not exceed 20 characters")
	}

	if input.Address != nil && len(*input.Address) > 255 {
		return errors.New("address must not exceed 255 characters")
	}

	return nil
}
