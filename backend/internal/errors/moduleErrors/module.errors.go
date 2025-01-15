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

func (h ModuleError) Error() string {
	if len(h.details) > 0 {
		return fmt.Sprintf("%s: %s", h.err.Error(), strings.Join(h.details, "; "))
	}
	return h.err.Error()
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
