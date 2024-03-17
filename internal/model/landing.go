package model

import "gorm.io/gorm"

type Landing struct {
	gorm.Model
	FlightID     uint
	ApproachType ApproachType
	Count        uint
	NightCount   uint
	DayCount     uint
	AirportCode  string
}
