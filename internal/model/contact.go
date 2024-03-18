package model

import "gorm.io/gorm"

// dany kontakt nalkezy do usera 1 to many
type Contact struct {
	gorm.Model
	UserID       uint
	AvatarURL    string
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
