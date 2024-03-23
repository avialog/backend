package model

import "gorm.io/gorm"

type Landing struct {
	gorm.Model
	FlightID     uint         `gorm:"required; not null; default:null"`
	ApproachType ApproachType `gorm:"required; not null; default:null"`
	Count        uint         `gorm:"required; not null; default:null"`
	NightCount   uint
	DayCount     uint
	AirportCode  string
}
