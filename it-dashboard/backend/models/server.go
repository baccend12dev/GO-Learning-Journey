package models

import (
	"time"
)

type Server struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255)"`
	IP        string `gorm:"type:varchar(255)"`
	OS        string `gorm:"type:varchar(255)"`
	Location  string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Server represents the structure of a server in the database.
type CreateServerRequest struct {
	Name     string `json:"name" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	OS       string `json:"os" binding:"required"`
	Location string `json:"location" binding:"required"`
}
