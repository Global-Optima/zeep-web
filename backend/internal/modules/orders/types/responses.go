package types

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500OrderCreate         = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, data.CreateOperation.ToString())
	Response500OrderPaymentSuccess = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, "payment", "success")
	Response500OrderPaymentFail    = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, "payment", "fail")
	Response400OrderCustomerName   = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "customerName")
	Response400Order               = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent)
	Response201Order               = localization.NewResponseKey(http.StatusCreated, data.OrderComponent)
	Response200OrderPaymentSuccess = localization.NewResponseKey(http.StatusOK, data.OrderComponent, "payment", "success")
	Response200OrderPaymentFail    = localization.NewResponseKey(http.StatusOK, data.OrderComponent, "payment", "fail")
	Response200OrderUpdate         = localization.NewResponseKey(http.StatusOK, data.OrderComponent, data.UpdateOperation.ToString())
	Response400InsufficientStock   = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "INSUFFICIENT_STOCK")
	Response400MultipleSelect      = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "MULTIPLE_SELECT")
)
