package moduleErrors

type ModuleErrorInterface interface {
	Error() string
}

type ModuleError struct {
	err error
}

func (h ModuleError) Error() string {
	return h.err.Error()
}

func NewModuleError(err error) *ModuleError {
	return &ModuleError{
		err: err,
	}
}
