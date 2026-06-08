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
	Server      Server `gorm:"foreignKey:ServerId"`
	Status      string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	Notes           []Note           `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` //delete the child to
	FeatureRequests []FeatureRequest `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"feature_requests,omitempty"`
	Documentations  []Documentation  `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"documentations,omitempty"`
	CreatedAt       time.Time
	UpdatedAt   time.Time
}

// System represents the structure of a system in the database.
// It includes fields such as ID, Name, Type, Links, ServerId, Status, Description, CreatedAt, and UpdatedAt.
type CreateSystemRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Links       string `json:"links" binding:"required"`
	ServerId    uint   `json:"server_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateSystemRequest represents the expected structure of the request body when creating a new system.
// It includes fields such as Name, Type, Links, ServerId, Status, and Description, all of which are required.
