package types

type SKUFilter struct {
	Name           *string
	Category       *string
	LowStock       *bool
	ExpirationFlag *bool
	IsActive       *bool
}
