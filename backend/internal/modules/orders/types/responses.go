package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"net/http"
)

var (
	Response500OrderCreate       = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, data.CreateOperation.ToString())
	Response400OrderCustomerName = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "customerName")
	Response400Order             = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent)
	Response201Order             = localization.NewResponseKey(http.StatusCreated, data.OrderComponent)
	Response200OrderUpdate       = localization.NewResponseKey(http.StatusOK, data.OrderComponent, data.UpdateOperation.ToString())
)
