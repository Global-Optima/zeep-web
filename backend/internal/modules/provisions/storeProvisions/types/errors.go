package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrStoreProvisionNotFound                 = moduleErrors.NewModuleError(errors.New("store provision not found"))
	ErrProvisionCompleted                     = moduleErrors.NewModuleError(errors.New("store provision has completed status"))
	ErrStoreProvisionIngredientMismatch       = moduleErrors.NewModuleError(errors.New("selected ingredient does not match the actual provision ingredients from central catalog"))
	ErrStoreProvisionDailyLimitReached        = moduleErrors.NewModuleError(errors.New("selected daily limit reached"))
	ErrInsufficientStoreProvision             = moduleErrors.NewModuleError(errors.New("insufficient store provision volume"))
	ErrInvalidStoreProvisionIngredientsVolume = moduleErrors.NewModuleError(errors.New("invalid store provision ingredients volume"))
)
