package data

type EmployeeWorkday struct {
	ID         uint   `gorm:"primaryKey"`
	Day        string `gorm:"size:15;not null"`
	StartAt    string `gorm:"type:time;not null"`
	EndAt      string `gorm:"type:time;not null"`
	EmployeeID uint   `gorm:"column:employee_id;not null"`
}
