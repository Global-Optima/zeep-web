package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrFailedToCheckRegionWarehouse = handlerErrors.NewHandlerError(errors.New("failed to check region warehouse"), http.StatusInternalServerError)
	ErrRegionWarehouseMismatch      = handlerErrors.NewHandlerError(errors.New("this warehouse is not assigned to the given region"), http.StatusForbidden)
)
