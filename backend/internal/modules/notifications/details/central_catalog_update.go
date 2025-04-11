package details

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type CentralCatalogUpdateDetails struct {
	BaseNotificationDetails
	Changes []CentralCatalogChange `json:"changes"`
}

type CentralCatalogChange struct {
	Key    string                 `json:"key"`
	Params map[string]interface{} `json:"params"`
}

func (c *CentralCatalogUpdateDetails) ToDetails() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CentralCatalogUpdateDetails) GetBaseDetails() *BaseNotificationDetails {
	return &c.BaseNotificationDetails
}

func BuildCentralCatalogUpdateDetails(facilityID uint, facilityName string, changes []CentralCatalogChange) (*CentralCatalogUpdateDetails, error) {
	if facilityID == 0 || facilityName == "" || changes == nil {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, updatedBy, and changes are required")
	}

	return &CentralCatalogUpdateDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		Changes: changes,
	}, nil
}

func BuildCentralCatalogUpdateChanges(changes []CentralCatalogChange) (localization.LocalizedMessage, error) {
	localizedMessages := localization.LocalizedMessage{}

	languages := map[string]*string{
		"en": &localizedMessages.En,
		"ru": &localizedMessages.Ru,
		"kk": &localizedMessages.Kk,
	}

	for lang, msg := range languages {
		var changeMessages []string

		for _, change := range changes {
			modifiedParams := make(map[string]interface{})
			for k, v := range change.Params {
				if strVal, ok := v.(string); ok && strVal == "" {
					emptyTranslated, err := localization.Translate("notification.emptyValue", nil)
					if err != nil {
						return localization.LocalizedMessage{}, fmt.Errorf("failed to translate empty value for language %s: %w", lang, err)
					}

					var localizedEmpty string
					switch lang {
					case "en":
						localizedEmpty = emptyTranslated.En
					case "ru":
						localizedEmpty = emptyTranslated.Ru
					case "kk":
						localizedEmpty = emptyTranslated.Kk
					default:
						localizedEmpty = emptyTranslated.Ru
					}
					modifiedParams[k] = localizedEmpty
				} else {
					modifiedParams[k] = v
				}
			}

			// Now translate using the (possibly) modified parameters.
			message, err := localization.Translate(change.Key, modifiedParams)
			if err != nil {
				return localization.LocalizedMessage{}, fmt.Errorf("failed to localize change %s for language %s: %w", change.Key, lang, err)
			}

			switch lang {
			case "en":
				changeMessages = append(changeMessages, message.En)
			case "ru":
				changeMessages = append(changeMessages, message.Ru)
			case "kk":
				changeMessages = append(changeMessages, message.Kk)
			}
		}

		*msg = strings.Join(changeMessages, "; ")
	}

	return localizedMessages, nil
}

func BuildCentralCatalogUpdateMessage(details *CentralCatalogUpdateDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	changesSummary, err := BuildCentralCatalogUpdateChanges(details.Changes)
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build changes summary: %w", err)
	}

	localizedMessages := localization.LocalizedMessage{}

	languages := map[string]*string{
		"en": &localizedMessages.En,
		"ru": &localizedMessages.Ru,
		"kk": &localizedMessages.Kk,
	}

	for lang, msg := range languages {
		key := localization.FormTranslationKey("notification", data.CENTRAL_CATALOG_UPDATE.ToString())
		translatedMessage, err := localization.Translate(key, map[string]interface{}{
			"FacilityName": details.FacilityName,
			"ID":           details.ID,
			"Changes":      changesSummaryField(lang, changesSummary),
		})
		if err != nil {
			return localization.LocalizedMessage{}, fmt.Errorf("failed to build %s message: %w", lang, err)
		}

		*msg = getMessageForLang(lang, *translatedMessage)
	}

	return localizedMessages, nil
}

func changesSummaryField(lang string, changes localization.LocalizedMessage) string {
	switch lang {
	case "en":
		return changes.En
	case "ru":
		return changes.Ru
	case "kk":
		return changes.Kk
	default:
		return changes.Ru
	}
}

func getMessageForLang(lang string, message localization.LocalizedMessage) string {
	switch lang {
	case "en":
		return message.En
	case "ru":
		return message.Ru
	case "kk":
		return message.Kk
	default:
		return message.Ru
	}
}
