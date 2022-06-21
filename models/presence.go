package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Presence Struct
type Presence struct {
	ID        string    `gorm:"type:varchar(100);column:id;primary_key:true"`
	UserID    string    `gorm:"type:varchar(100);column:user_id"`
	PhotoID   string    `gorm:"type:varchar(100);column:photo_id"`
	CheckIN   time.Time `gorm:"column:checkin"`
	CheckOut  time.Time `gorm:"column:checkout"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (c Presence) TableName() string {
	return "presence"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Presence) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
