package data

import "time"

type EmployeeAudit struct {
	ID          uint       `gorm:"primaryKey"`
	StartWorkAt *time.Time `gorm:"type:timestamp"`
	EndWorkAt   *time.Time `gorm:"type:timestamp"`
	EmployeeID  uint       `gorm:"column:employee_id;not null"`
}
