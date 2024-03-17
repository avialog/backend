package model

import "gorm.io/gorm"

type Aircraft struct {
	gorm.Model
	UserID             int64
	RegistrationNumber string
	AircraftModel      string
	Remarks            string // Notes
	ImageURL           string
}
