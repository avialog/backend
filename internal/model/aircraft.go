package model

import "gorm.io/gorm"

type Aircraft struct {
	gorm.Model
	UserID             uint `gorm:"required; not null; default:null"`
	User               User
	RegistrationNumber string `gorm:"required; not null; default:null"`
	AircraftModel      string
	Remarks            string
	ImageURL           string
	Flights            []Flight `gorm:"foreignKey:AircraftID"`
}
