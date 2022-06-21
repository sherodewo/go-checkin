package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Photo Struct
type Photo struct {
	ID        string    `gorm:"type:varchar(100);column:id;primary_key:true"`
	UserID    string    `gorm:"type:varchar(100);column:user_id"`
	Path      string    `gorm:"type:varchar(255);column:path"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (c Photo) TableName() string {
	return "photo"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
