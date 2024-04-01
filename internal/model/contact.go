package model

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID       string `gorm:"required; not null; default:null" validate:"required"`
	User         User   `validate:"-"`
	AvatarURL    string
	FirstName    string `gorm:"required; not null; default:null" validate:"required"`
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
