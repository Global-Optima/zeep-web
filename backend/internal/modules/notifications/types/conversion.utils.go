package types

import (
	"encoding/json"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertToNotificationDetails(details *data.EmployeeNotification) (*data.ExtendedDetails, error) {
	var extendedDetails data.ExtendedDetails
	err := json.Unmarshal(details.Details, &extendedDetails)
	if err != nil {
		return nil, err
	}
	return &extendedDetails, nil
}
