package model

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID       uint
	User         User
	AvatarURL    string
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
