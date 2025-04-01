package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var Response500AuditGet = localization.NewResponseKey(500, data.AuditComponent, data.GetOperation.ToString())
