package handlerErrors

type HandlerErrorInterface interface {
	Error() string
	Status() int
}

type HandlerError struct {
	err    error
	status int
}

func (h HandlerError) Error() string {
	return h.err.Error()
}

func (h HandlerError) Status() int {
	return h.status
}

func NewHandlerError(err error, status int) *HandlerError {
	return &HandlerError{
		err:    err,
		status: status,
	}
}
