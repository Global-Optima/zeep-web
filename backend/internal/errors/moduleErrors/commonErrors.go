package moduleErrors

import "github.com/pkg/errors"

var (
	ErrNotFound      = NewModuleError(errors.New("resource not found"))
	ErrValidation    = NewModuleError(errors.New("validation error"))
	ErrAlreadyExists = NewModuleError(errors.New("object already exists"))
)
