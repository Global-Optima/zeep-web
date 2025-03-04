package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	// 500 Internal Server Error responses
	Response500SupplierCreate          = localization.NewResponseKey(500, data.SupplierComponent, data.CreateOperation.ToString())
	Response500SupplierGet             = localization.NewResponseKey(500, data.SupplierComponent, data.GetOperation.ToString())
	Response500SupplierUpdate          = localization.NewResponseKey(500, data.SupplierComponent, data.UpdateOperation.ToString())
	Response500SupplierDelete          = localization.NewResponseKey(500, data.SupplierComponent, data.DeleteOperation.ToString())
	Response500SupplierGetMaterials    = localization.NewResponseKey(500, data.SupplierComponent, "SUPPLIER_MATERIALS")
	Response500SupplierUpsertMaterials = localization.NewResponseKey(500, data.SupplierComponent, "UPSERT_MATERIALS")
	// 404 Not Found response
	Response404Supplier = localization.NewResponseKey(404, data.SupplierComponent)
	// 400 Bad Request responses
	Response400Supplier = localization.NewResponseKey(400, data.SupplierComponent, data.CreateOperation.ToString())
	// 201 Created response
	Response201Supplier = localization.NewResponseKey(201, data.SupplierComponent)
	// 200 Success responses
	Response200SupplierUpdate         = localization.NewResponseKey(200, data.SupplierComponent, data.UpdateOperation.ToString())
	Response200SupplierDelete         = localization.NewResponseKey(200, data.SupplierComponent, data.DeleteOperation.ToString())
	Response200SupplierUpsertMaterial = localization.NewResponseKey(200, data.SupplierComponent, "UPSERT_MATERIALS")
)
