package translations

type FieldLocale struct {
	En string `json:"en" binding:"omitempty"`
	Ru string `json:"ru" binding:"omitempty"`
	Kk string `json:"kk" binding:"omitempty"`
}
