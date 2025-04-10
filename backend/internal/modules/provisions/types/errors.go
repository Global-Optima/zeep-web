package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrProvisionNotFound   = moduleErrors.NewModuleError(errors.New("provision not found"))
	ErrProvisionUniqueName = moduleErrors.NewModuleError(errors.New("provision with the given name already exists"))
	ErrProvisionIsInUse    = moduleErrors.NewModuleError(errors.New("provision is in use"))
)
