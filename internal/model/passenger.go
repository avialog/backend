package model

type Passenger struct {
	ID           int64
	FlightID     int64
	Role         Role
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
