package data

type Store struct {
	BaseEntity
	Name              string             `gorm:"size:255;not null" sort:"name"`
	FacilityAddressID uint               `gorm:"index;not null"`
	FacilityAddress   FacilityAddress    `gorm:"foreignKey:FacilityAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsFranchise       bool               `gorm:"default:false" sort:"isFranchise"`
	AdminID           *uint              `gorm:"index;not null"`
	ContactPhone      string             `gorm:"size:20"`
	ContactEmail      string             `gorm:"size:255"`
	StoreHours        string             `gorm:"size:255"`
	Additives         []StoreAdditive    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	ProductSizes      []StoreProductSize `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Products          []StoreProduct     `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
}

type StoreAdditive struct {
	BaseEntity
	AdditiveID uint     `gorm:"index;not null"`
	StoreID    uint     `gorm:"index;not null"`
	Price      float64  `gorm:"type:decimal(10,2);default:0"`
	Store      Store    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Additive   Additive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
}

type StoreProductSize struct {
	BaseEntity
	ProductSizeID uint        `gorm:"index;not null"`
	StoreID       uint        `gorm:"index;not null"`
	Price         float64     `gorm:"type:decimal(10,2);default:0"`
	Store         Store       `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	ProductSize   ProductSize `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type StoreProduct struct {
	BaseEntity
	ProductID   uint    `gorm:"index;not null"`
	StoreID     uint    `gorm:"index;not null"`
	IsAvailable bool    `gorm:"default:true"`
	Store       Store   `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type FacilityAddress struct {
	BaseEntity
	Address   string   `gorm:"size:255;not null"`
	Longitude *float64 `gorm:"type:decimal(9,6)"`
	Latitude  *float64 `gorm:"type:decimal(9,6)"`
	Stores    []Store  `gorm:"foreignKey:FacilityAddressID"`
}
