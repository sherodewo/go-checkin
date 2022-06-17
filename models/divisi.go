package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Divisi Struct
type Divisi struct {
	ID          string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	Name        string    `gorm:"type:varchar(50);column:name"`
	Description string    `gorm:"type:varchar(255);column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (c Divisi) TableName() string {
	return "divisi"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Divisi) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
