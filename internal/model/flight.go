package model

import "time"

type Flight struct {
	ID                  int64
	UserID              int64
	AircraftID          int64
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
