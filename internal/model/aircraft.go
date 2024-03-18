package model

import "gorm.io/gorm"

type Aircraft struct {
	gorm.Model
	UserID             uint
	User               User
	RegistrationNumber string
	AircraftModel      string
	Remarks            string
	ImageURL           string
	Flights            []Flight `gorm:"foreignKey:AircraftID"`
}
