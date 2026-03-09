package models

import (
	"time"
)

type Note struct {
	ID        uint `gorm:"primaryKey"`
	SystemId  uint
	Title     string `gorm:"type:varchar(255)"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Note represents the structure of a note in the database.
// It includes fields such as ID, SystemId, Title, Content, CreatedAt, and UpdatedAt.
type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
