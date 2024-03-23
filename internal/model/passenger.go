package model

import "gorm.io/gorm"

type Passenger struct {
	gorm.Model
	FlightID     uint `gorm:"required; not null; default:null"`
	Flight       Flight
	Role         Role   `gorm:"required; not null; default:null"`
	FirstName    string `gorm:"required; not null; default:null"`
	LastName     string `gorm:"required; not null; default:null"`
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
