package model

import (
	"gorm.io/gorm"
	"time"
)

type Flight struct {
	gorm.Model
	UserID              uint
	AircraftID          uint
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
