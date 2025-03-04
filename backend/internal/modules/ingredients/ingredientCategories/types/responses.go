package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response400IngredientCategory            = localization.NewResponseKey(400, data.IngredientCategoryComponent)
	Response500IngredientCategoryGet         = localization.NewResponseKey(500, data.IngredientCategoryComponent, data.GetOperation.ToString())
	Response500IngredientCategoryCreate      = localization.NewResponseKey(500, data.IngredientCategoryComponent, data.CreateOperation.ToString())
	Response500IngredientCategoryUpdate      = localization.NewResponseKey(500, data.IngredientCategoryComponent, data.UpdateOperation.ToString())
	Response500IngredientCategoryDelete      = localization.NewResponseKey(500, data.IngredientCategoryComponent, data.DeleteOperation.ToString())
	Response409IngredientCategoryDeleteInUse = localization.NewResponseKey(409, data.IngredientCategoryComponent, data.DeleteOperation.ToString(), "in-use")
	Response200IngredientCategoryUpdate      = localization.NewResponseKey(200, data.IngredientCategoryComponent, data.UpdateOperation.ToString())
	Response200IngredientCategoryDelete      = localization.NewResponseKey(200, data.IngredientCategoryComponent, data.DeleteOperation.ToString())
	Response201IngredientCategory            = localization.NewResponseKey(201, data.IngredientCategoryComponent)
)
