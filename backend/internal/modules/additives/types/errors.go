package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrAdditiveAlreadyExists = moduleErrors.NewModuleError(errors.New("additive already exists"))
)
