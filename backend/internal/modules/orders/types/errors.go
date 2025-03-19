package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrOrderNotFound             = moduleErrors.NewModuleError(errors.New("order not found"))
	ErrInappropriateOrderStatus  = moduleErrors.NewModuleError(errors.New("inappropriate order status"))
	ErrInsufficientStock         = moduleErrors.NewModuleError(errors.New("insufficient stock to fulfill the order"))
	ErrMultipleSelect            = moduleErrors.NewModuleError(errors.New("multiple select on this additive category is not allowed"))
	ErrInvalidCustomerNameCensor = moduleErrors.NewModuleError(errors.New("invalid customer name"))
)
