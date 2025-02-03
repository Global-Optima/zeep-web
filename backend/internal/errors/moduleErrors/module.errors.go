package moduleErrors

import (
	"fmt"
	"strings"
)

type ModuleErrorInterface interface {
	Error() string
	WithDetails(details ...string) ModuleErrorInterface
	Details() []string
}

type ModuleError struct {
	err     error
	details []string
}

func (m ModuleError) Error() string {
	if len(m.details) > 0 {
		return fmt.Sprintf("%s: %s", m.err.Error(), strings.Join(m.details, "; "))
	}
	return m.err.Error()
}

func (m ModuleError) Details() []string {
	return m.details
}

func (m ModuleError) WithDetails(details ...string) ModuleErrorInterface {
	return &ModuleError{
		err:     m.err,
		details: append(m.details, details...),
	}
}

func NewModuleError(err error) *ModuleError {
	return &ModuleError{
		err: err,
	}
}
