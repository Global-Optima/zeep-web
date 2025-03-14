package moduleErrors

import (
	"fmt"
	"strings"
)

type ModuleErrorInterface interface {
	Error() string
	WithDetails(reason string, details ...string) ModuleErrorInterface
	GetDetails() string
}

type ModuleError struct {
	err     error
	reason  string
	details string
}

func (m ModuleError) Error() string {
	if len(m.details) > 0 {
		return fmt.Sprintf("%s: %s", m.err.Error(), m.details)
	}
	return m.err.Error()
}

func (m ModuleError) WithDetails(reason string, details ...string) ModuleErrorInterface {
	return &ModuleError{
		err:     m.err,
		reason:  reason,
		details: strings.Join(details, "; "),
	}
}

func (m ModuleError) GetDetails() string {
	return m.details
}

func NewModuleError(err error) *ModuleError {
	return &ModuleError{
		err: err,
	}
}
