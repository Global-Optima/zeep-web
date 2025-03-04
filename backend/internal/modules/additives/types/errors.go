package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrAdditiveCategoryIsInUse = moduleErrors.NewModuleError(errors.New("additive category is in use"))
	ErrAdditiveAlreadyExists   = moduleErrors.NewModuleError(errors.New("additive already exists"))
	ErrAdditiveNotFound        = moduleErrors.NewModuleError(errors.New("additive not found"))

	ErrAdditiveCategoryNotFound = moduleErrors.NewModuleError(errors.New("additive category not found"))
)
