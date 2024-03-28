package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
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
