package types

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/pkg/errors"
)

var (
	ErrFailedToCheckFranchiseeStore = handlerErrors.NewHandlerError(errors.New("failed to check franchise store"), http.StatusInternalServerError)
	ErrFranchiseeStoreMismatch      = handlerErrors.NewHandlerError(errors.New("this store is not assigned to the given franchisee"), http.StatusForbidden)
)
