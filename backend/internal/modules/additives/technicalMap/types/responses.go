package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var Response500TechnicalMapGet = localization.NewResponseKey(500, data.TechnicalMapComponent, data.GetOperation.ToString())
