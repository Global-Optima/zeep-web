package data

type RecipeStep struct {
	ID          uint   `gorm:"primaryKey"`
	ProductID   uint   `gorm:"column:product_id;not null"`
	Step        int    `gorm:"not null"`
	Name        string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	ImageURL    string `gorm:"size:2048"`
}
