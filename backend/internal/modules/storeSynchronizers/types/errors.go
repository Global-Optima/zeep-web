package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrStoreAlreadySynchronized = moduleErrors.NewModuleError(errors.New("store is already synchronized"))
)
