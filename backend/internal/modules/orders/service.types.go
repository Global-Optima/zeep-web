package orders

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditivesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
)

const (
	OrderPaymentFailure = "order-payment-failure"
)

type subordersQuantities struct {
	storeProductSizesQty map[uint]uint
	storeAdditivesQty    map[uint]uint
}

type preparedData struct {
	storeProductSizeIDs       []uint
	suborderStoreAdditivesCtx *suborderAdditivesContext
	suborderQuantities        *subordersQuantities
}

type storeProductSizeValidationResults struct {
	storeProductSizesList []data.StoreProductSize
	prices                map[uint]float64
	names                 map[uint]string
}

type storeAdditiveValidationResults struct {
	storeAdditivesList []data.StoreAdditive
	prices             map[uint]float64
	names              map[uint]string
}

type suborderAdditivesContext struct {
	storeAddKeys []storeAdditivesTypes.StorePStoAdditiveKey
	storeAddIDs  []uint
}

type subordersContext struct {
	subordersQuantities   *subordersQuantities
	storeProductSizesList []data.StoreProductSize
	storeAdditivesList    []data.StoreAdditive
}

type productSizeToAdditiveKey struct {
	productSizeID uint
	additiveID    uint
}

type orderValidationResults struct {
	productPrices  map[uint]float64
	productNames   map[uint]string
	additivePrices map[uint]float64
	additiveNames  map[uint]string
	subordersCtx   *subordersContext
}
