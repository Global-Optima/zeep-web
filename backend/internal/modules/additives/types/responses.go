package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500Additive       = localization.NewResponseKey(500, data.AdditiveComponent)
	Response400Additive       = localization.NewResponseKey(400, data.AdditiveComponent)
	Response200AdditiveUpdate = localization.NewResponseKey(200, data.AdditiveComponent, data.UpdateOperation.ToString())
	Response200AdditiveDelete = localization.NewResponseKey(200, data.AdditiveComponent, data.DeleteOperation.ToString())
	Response201Additive       = localization.NewResponseKey(201, data.AdditiveComponent)

	Response500AdditiveCategory       = localization.NewResponseKey(500, data.AdditiveCategoryComponent)
	Response400AdditiveCategory       = localization.NewResponseKey(400, data.AdditiveCategoryComponent)
	Response200AdditiveCategoryUpdate = localization.NewResponseKey(200, data.AdditiveCategoryComponent, data.UpdateOperation.ToString())
	Response200AdditiveCategoryDelete = localization.NewResponseKey(200, data.AdditiveCategoryComponent, data.DeleteOperation.ToString())
	Response201AdditiveCategory       = localization.NewResponseKey(201, data.AdditiveCategoryComponent)
)
