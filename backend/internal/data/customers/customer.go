package data

type Customer struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:255;not null"`
	Password   string `gorm:"size:255;not null"`
	Phone      string `gorm:"size:15;unique"`
	IsVerified bool   `gorm:"default:false"`
	IsBanned   bool   `gorm:"default:false"`
}
