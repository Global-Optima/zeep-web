package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
)

func ValidateCustomer(input CustomerRegisterDTO) error {
	if input.Name == "" {
		return errors.New("customer name cannot be empty")
	}

	if !utils.IsValidPhone(input.Phone) {
		return errors.New("invalid email format")
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	return nil
}
