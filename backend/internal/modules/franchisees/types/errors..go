package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrFailedToCheckFranchiseeStore = handlerErrors.NewHandlerError(errors.New("failed to check franchise store"), http.StatusInternalServerError)
	ErrFranchiseeStoreMismatch      = handlerErrors.NewHandlerError(errors.New("this store is not assigned to the given franchisee"), http.StatusForbidden)
)
