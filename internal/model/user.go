package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	AvatarURL    string
	SignatureURL string
	Country      Country
	Phone        string
	Street       string
	City         string
	Company      string
	Timezone     string
	Contacts     []Contact  `gorm:"foreignKey:UserID"`
	Aircraft     []Aircraft `gorm:"foreignKey:UserID"`
	Flights      []Flight   `gorm:"foreignKey:UserID"`
}
