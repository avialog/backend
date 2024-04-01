package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	FirstName    string
	LastName     string
	Email        string `gorm:"required; not null; default:null"`
	AvatarURL    string
	SignatureURL string
	Country      Country
	Phone        string
	Street       string
	City         string
	Company      string
	Timezone     string
	Contacts     []Contact  `gorm:"foreignKey:UserID" validate:"-"`
	Aircraft     []Aircraft `gorm:"foreignKey:UserID" validate:"-"`
	Flights      []Flight   `gorm:"foreignKey:UserID" validate:"-"`
}
