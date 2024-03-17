package model

import (
	"gorm.io/gorm"
	"time"
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
	Timezone     time.Location
}
