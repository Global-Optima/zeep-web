package handlerErrors

import (
	"fmt"
	"strings"
)

type HandlerErrorInterface interface {
	Error() string
	Status() int
	WithDetails(details ...string) HandlerErrorInterface
	Details() []string
}

type HandlerError struct {
	err     error
	status  int
	details []string
}

func (h HandlerError) Error() string {
	if len(h.details) > 0 {
		return fmt.Sprintf("%s: %s", h.err.Error(), strings.Join(h.details, "; "))
	}
	return h.err.Error()
}

func (h HandlerError) Status() int {
	return h.status
}

func (h HandlerError) Details() []string {
	return h.details
}

func (h HandlerError) WithDetails(details ...string) HandlerErrorInterface {
	return &HandlerError{
		err:     h.err,
		status:  h.status,
		details: append(h.details, details...),
	}
}

func NewHandlerError(err error, status int) *HandlerError {
	return &HandlerError{
		err:    err,
		status: status,
	}
}
