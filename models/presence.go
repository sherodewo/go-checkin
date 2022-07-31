package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Presence Struct
type Presence struct {
	ID           string    `gorm:"column:id;primary_key:true"`
	UserID       string    `gorm:"column:user_id"`
	Checkin      time.Time `gorm:"column:checkin"`
	Checkout     time.Time `gorm:"column:checkout"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	PhotoID      string    `gorm:"column:photo_id"`
	LocationName string    `gorm:"column:location_name"`
	Location     string    `gorm:"column:location"`
}

type PresenceDatatable struct {
	Name         string    `gorm:"type:varchar(100);column:name" json:"name"`
	ID           string    `gorm:"column:id;primary_key:true"`
	UserID       string    `gorm:"column:user_id"`
	Checkin      time.Time `gorm:"column:checkin"`
	Checkout     time.Time `gorm:"column:checkout"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	PhotoID      string    `gorm:"column:photo_id"`
	LocationName string    `gorm:"column:location_name"`
	Location     string    `gorm:"column:location"`
}

func (c Presence) TableName() string {
	return "presence"
}

func (c PresenceDatatable) TableName() string {
	return "presence"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Presence) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
