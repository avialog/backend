package model

import "gorm.io/gorm"

// londowanie nalezy do lotu
type Landing struct {
	gorm.Model
	FlightID     uint
	ApproachType ApproachType
	Count        uint
	NightCount   uint
	DayCount     uint
	AirportCode  string
}
