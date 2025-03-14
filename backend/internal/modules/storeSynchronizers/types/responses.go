package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

const SYNCHRONIZATION_KEY = "synchronization"

var (
	Response200StoreSynchronization      = localization.NewResponseKey(200, data.StoreComponent, SYNCHRONIZATION_KEY)
	Response500StoreSynchronization      = localization.NewResponseKey(500, data.StoreComponent, SYNCHRONIZATION_KEY)
	Response500StoreSynchronizationCheck = localization.NewResponseKey(500, data.StoreComponent, SYNCHRONIZATION_KEY, "check")
)
