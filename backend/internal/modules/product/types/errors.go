package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var ErrProductAlreadyExists = moduleErrors.NewModuleError(errors.New("product already exists"))
