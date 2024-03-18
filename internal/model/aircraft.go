package model

import "gorm.io/gorm"

// wielke userow moze miec ten sam samolot
type Aircraft struct {
	gorm.Model
	UserID             int64 //obcy
	RegistrationNumber string
	AircraftModel      string
	Remarks            string // Notes
	ImageURL           string
}
