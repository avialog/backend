package model

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID       uint `gorm:"required; not null; default:null"`
	User         User
	AvatarURL    string
	FirstName    string `gorm:"required; not null; default:null"`
	LastName     string `gorm:"required; not null; default:null"`
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
