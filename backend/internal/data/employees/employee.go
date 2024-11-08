package data

type Employee struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	Phone    string `gorm:"size:15;unique"`
	Email    string `gorm:"size:255;unique"`
	RoleID   *uint  `gorm:"column:role_id"`
	StoreID  *uint  `gorm:"column:store_id"`
	IsActive bool   `gorm:"default:true"`
}
