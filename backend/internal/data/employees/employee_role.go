package data

type EmployeeRole struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;not null"`
}
