package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedEmployeeType = moduleErrors.NewModuleError(errors.New("unsupported employee type"))
)
