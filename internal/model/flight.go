package model

import (
	"gorm.io/gorm"
	"time"
)

type Flight struct {
	gorm.Model
	UserID              uint `gorm:"required; not null; default:null"`
	User                User
	AircraftID          uint `gorm:"required; not null; default:null"`
	Aircraft            Aircraft
	Passengers          []Passenger `gorm:"foreignKey:FlightID"`
	Landings            []Landing   `gorm:"foreignKey:FlightID"`
	TakeoffTime         time.Time   `gorm:"required; not null; default:null"`
	TakeoffAirportCode  string      `gorm:"required; not null; default:null"`
	LandingTime         time.Time   `gorm:"required; not null; default:null"`
	LandingAirportCode  string      `gorm:"required; not null; default:null"`
	Style               Style       `gorm:"required; not null; default:null"`
	Remarks             string
	PersonalRemarks     string
	TotalBlockTime      time.Duration
	PilotInCommandTime  time.Duration
	SecondInCommandTime time.Duration
	DualReceivedTime    time.Duration
	DualGivenTime       time.Duration
	MultiPilotTime      time.Duration
	NightTime           time.Duration
	IFRTime             time.Duration
	IFRActualTime       time.Duration
	IFRSimulatedTime    time.Duration
	CrossCountryTime    time.Duration
	SimulatorTime       time.Duration
	SignatureURL        string
}
