package localization

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
)

var ComponentNameTranslations = map[Locale]map[data.ComponentName]string{
	English: {
		data.ProductComponent:             "Product",
		data.StoreProductComponent:        "Store Product",
		data.EmployeeComponent:            "Employee",
		data.AdditiveComponent:            "Additive",
		data.StoreAdditiveComponent:       "Store Additive",
		data.ProductSizeComponent:         "Product Size",
		data.StoreProductSizeComponent:    "Store Product Size",
		data.RecipeStepsComponent:         "Recipe Steps",
		data.StoreComponent:               "Store",
		data.WarehouseComponent:           "Warehouse",
		data.StoreWarehouseStockComponent: "Store Warehouse Stock",
		data.IngredientComponent:          "Ingredient",
	},
	Russian: {
		data.ProductComponent:             "Продукт",
		data.StoreProductComponent:        "Товар в магазине",
		data.EmployeeComponent:            "Сотрудник",
		data.AdditiveComponent:            "Добавка",
		data.StoreAdditiveComponent:       "Добавка в магазине",
		data.ProductSizeComponent:         "Размер продукта",
		data.StoreProductSizeComponent:    "Размер товара в магазине",
		data.RecipeStepsComponent:         "Этапы рецепта",
		data.StoreComponent:               "Магазин",
		data.WarehouseComponent:           "Склад",
		data.StoreWarehouseStockComponent: "Складской запас магазина",
		data.IngredientComponent:          "Ингредиент",
	},
	Kazakh: {
		data.ProductComponent:             "Өнім",
		data.StoreProductComponent:        "Дүкен өнімі",
		data.EmployeeComponent:            "Қызметкер",
		data.AdditiveComponent:            "Қосымша",
		data.StoreAdditiveComponent:       "Дүкен қосымшасы",
		data.ProductSizeComponent:         "Өнім мөлшері",
		data.StoreProductSizeComponent:    "Дүкендегі өнім мөлшері",
		data.RecipeStepsComponent:         "Рецепт қадамдары",
		data.StoreComponent:               "Дүкен",
		data.WarehouseComponent:           "Қойма",
		data.StoreWarehouseStockComponent: "Дүкеннің қойма қоры",
		data.IngredientComponent:          "Ингредиент",
	},
}

func GetLocalizedComponentName(locale Locale, componentName data.ComponentName) string {
	zapLogger := logger.GetZapSugaredLogger()

	translations, ok := ComponentNameTranslations[locale]
	if !ok {
		translations = ComponentNameTranslations[DEFAULT_LOCALE]
	}

	localizedName, ok := translations[componentName]
	if !ok {
		zapLogger.Warnf("Translation not found for component name: %s", componentName)
		return string(componentName)
	}

	return localizedName
}

var OperationTypeTranslations = map[Locale]map[data.OperationType]string{
	English: {
		data.CreateOperation:         "Create",
		data.UpdateOperation:         "Update",
		data.DeleteOperation:         "Delete",
		data.CreateMultipleOperation: "Create Multiple",
	},
	Russian: {
		data.CreateOperation:         "Создать",
		data.UpdateOperation:         "Обновить",
		data.DeleteOperation:         "Удалить",
		data.CreateMultipleOperation: "Создать несколько",
	},
	Kazakh: {
		data.CreateOperation:         "Жасау",
		data.UpdateOperation:         "Жаңарту",
		data.DeleteOperation:         "Жою",
		data.CreateMultipleOperation: "Бірнеше жасау",
	},
}

func GetLocalizedOperationType(locale Locale, operationType data.OperationType) string {
	zapLogger := logger.GetZapSugaredLogger()

	translations, ok := OperationTypeTranslations[locale]
	if !ok {
		translations = OperationTypeTranslations[DEFAULT_LOCALE]
	}

	localizedName, ok := translations[operationType]
	if !ok {
		zapLogger.Warnf("Translations not found for operation type: %s", operationType)
		return string(operationType)
	}

	return localizedName
}
