package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500AdditiveCreate      = localization.NewResponseKey(500, data.AdditiveComponent, data.CreateOperation.ToString())
	Response500AdditiveUpdate      = localization.NewResponseKey(500, data.AdditiveComponent, data.UpdateOperation.ToString())
	Response500AdditiveGet         = localization.NewResponseKey(500, data.AdditiveComponent, data.GetOperation.ToString())
	Response500AdditiveDelete      = localization.NewResponseKey(500, data.AdditiveComponent, data.DeleteOperation.ToString())
	Response409AdditiveDeleteInUse = localization.NewResponseKey(409, data.AdditiveComponent, data.DeleteOperation.ToString(), "IN_USE")
	Response404Additive            = localization.NewResponseKey(404, data.AdditiveComponent)
	Response400Additive            = localization.NewResponseKey(400, data.AdditiveComponent)
	Response200AdditiveUpdate      = localization.NewResponseKey(200, data.AdditiveComponent, data.UpdateOperation.ToString())
	Response200AdditiveDelete      = localization.NewResponseKey(200, data.AdditiveComponent, data.DeleteOperation.ToString())
	Response201Additive            = localization.NewResponseKey(201, data.AdditiveComponent)

	Response500AdditiveCategoryCreate      = localization.NewResponseKey(500, data.AdditiveCategoryComponent, data.CreateOperation.ToString())
	Response500AdditiveCategoryUpdate      = localization.NewResponseKey(500, data.AdditiveCategoryComponent, data.UpdateOperation.ToString())
	Response500AdditiveCategoryGet         = localization.NewResponseKey(500, data.AdditiveCategoryComponent, data.GetOperation.ToString())
	Response500AdditiveCategoryDelete      = localization.NewResponseKey(500, data.AdditiveCategoryComponent, data.DeleteOperation.ToString())
	Response409AdditiveCategoryDeleteInUse = localization.NewResponseKey(409, data.AdditiveCategoryComponent, data.DeleteOperation.ToString(), "in-use")
	Response404AdditiveCategory            = localization.NewResponseKey(404, data.AdditiveCategoryComponent)
	Response400AdditiveCategory            = localization.NewResponseKey(400, data.AdditiveCategoryComponent)
	Response200AdditiveCategoryUpdate      = localization.NewResponseKey(200, data.AdditiveCategoryComponent, data.UpdateOperation.ToString())
	Response200AdditiveCategoryDelete      = localization.NewResponseKey(200, data.AdditiveCategoryComponent, data.DeleteOperation.ToString())
	Response201AdditiveCategory            = localization.NewResponseKey(201, data.AdditiveCategoryComponent)

	Response500AdditiveTranslationsUpdate = localization.NewResponseKey(500, data.AdditiveComponent, data.UpdateOperation.ToString(), "translations")
	Response200AdditiveTranslationsUpdate = localization.NewResponseKey(200, data.AdditiveComponent, data.UpdateOperation.ToString(), "translations")
)
