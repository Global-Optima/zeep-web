package data

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
}
