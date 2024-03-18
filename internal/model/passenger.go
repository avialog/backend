package model

import "gorm.io/gorm"

// pasazer nalezy do lotu
type Passenger struct {
	gorm.Model
	FlightID     uint
	Flight       Flight
	Role         Role
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
