package model

import "gorm.io/gorm"

type Landing struct {
	gorm.Model
	FlightID     uint         `gorm:"required; not null; default:null" validate:"required"`
	ApproachType ApproachType `gorm:"required; not null; default:null" validate:"required,approach_type"`
	Count        uint
	NightCount   uint
	DayCount     uint
	AirportCode  string
}
