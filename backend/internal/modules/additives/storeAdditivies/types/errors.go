package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrStoreAdditiveAlreadyExists      = moduleErrors.NewModuleError(errors.New("Store additive already exists"))
	ErrStoreAdditiveCategoriesNotFound = moduleErrors.NewModuleError(errors.New("Store additive categories not found"))
)
