package types

type ProductFilterDao struct {
	StoreID     uint
	CategoryID  *uint
	SearchQuery string
	Limit       int
	Offset      int
}
