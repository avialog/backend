package model

type Landing struct {
	ID           int64
	FlightID     int64
	ApproachType ApproachType
	Count        uint
	NightCount   uint
	DayCount     uint
	AirportCode  string
}
