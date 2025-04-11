package types

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500OrderCreate         = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, data.CreateOperation.ToString())
	Response500OrderPaymentSuccess = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, "payment", "success")
	Response500SuborderNextStatus  = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, "next_status")
	Response500OrderPaymentFail    = localization.NewResponseKey(http.StatusInternalServerError, data.OrderComponent, "payment", "fail")
	Response404Order               = localization.NewResponseKey(http.StatusNotFound, data.OrderComponent)
	Response400OrderCustomerName   = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "CUSTOMER_NAME")
	Response400Order               = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent)
	Response201Order               = localization.NewResponseKey(http.StatusCreated, data.OrderComponent)
	Response200OrderPaymentSuccess = localization.NewResponseKey(http.StatusOK, data.OrderComponent, "payment", "success")
	Response200OrderPaymentFail    = localization.NewResponseKey(http.StatusOK, data.OrderComponent, "payment", "fail")
	Response200OrderUpdate         = localization.NewResponseKey(http.StatusOK, data.OrderComponent, data.UpdateOperation.ToString())
	Response409InsufficientStock   = localization.NewResponseKey(http.StatusConflict, data.OrderComponent, "INSUFFICIENT_STOCK")
	Response400MultipleSelect      = localization.NewResponseKey(http.StatusBadRequest, data.OrderComponent, "MULTIPLE_SELECT")
)
