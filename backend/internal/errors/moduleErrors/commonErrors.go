package moduleErrors

import "github.com/pkg/errors"

var (
	ErrNotFound = NewModuleError(errors.New("resource not found"))
)
