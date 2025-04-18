package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response400Ingredient            = localization.NewResponseKey(400, data.IngredientComponent)
	Response500IngredientCreate      = localization.NewResponseKey(500, data.IngredientComponent, data.CreateOperation.ToString())
	Response500IngredientUpdate      = localization.NewResponseKey(500, data.IngredientComponent, data.UpdateOperation.ToString())
	Response500IngredientDelete      = localization.NewResponseKey(500, data.IngredientComponent, data.DeleteOperation.ToString())
	Response404Ingredient            = localization.NewResponseKey(404, data.IngredientComponent)
	Response409IngredientDeleteInUse = localization.NewResponseKey(409, data.IngredientComponent, data.DeleteOperation.ToString(), "IN_USE")
	Response201IngredientCreate      = localization.NewResponseKey(201, data.IngredientComponent)
	Response200IngredientUpdate      = localization.NewResponseKey(200, data.IngredientComponent, data.UpdateOperation.ToString())
	Response200IngredientDelete      = localization.NewResponseKey(200, data.IngredientComponent, data.DeleteOperation.ToString())

	Response200IngredientTranslationsUpdate = localization.NewResponseKey(200, data.IngredientComponent, data.UpdateOperation.ToString(), "TRANSLATIONS")
	Response500IngredientTranslationsUpdate = localization.NewResponseKey(500, data.IngredientComponent, data.UpdateOperation.ToString(), "TRANSLATIONS")
)
