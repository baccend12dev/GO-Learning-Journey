package models

import (
	"time"
)

type System struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255)"`
	Type        string `gorm:"type:varchar(255)"`
	Links       string `gorm:"type:varchar(255)"`
	ServerId    uint
	Status      string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
