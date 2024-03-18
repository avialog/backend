package model

import (
	"gorm.io/gorm"
	"time"
)

// lot nalezy do usera
type Flight struct {
	gorm.Model
	UserID              uint
	User                User
	AircraftID          uint
	Aircraft            Aircraft
	Passengers          []Passenger `gorm:"foreignKey:FlightID"`
	Landings            []Landing   `gorm:"foreignKey:FlightID"`
	TakeoffTime         time.Time
	TakeoffAirportCode  string
	LandingTime         time.Time
	LandingAirportCode  string
	Style               Style
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
