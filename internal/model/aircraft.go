package model

import "gorm.io/gorm"

type Aircraft struct {
	gorm.Model
	UserID             string `gorm:"required; not null; default:null" validate:"required"`
	User               User   `validate:"-"`
	RegistrationNumber string `gorm:"required; not null; default:null" validate:"required"`
	AircraftModel      string `gorm:"required; not null; default:null" validate:"required"`
	Remarks            *string
	ImageURL           *string
	Flights            []Flight `gorm:"foreignKey:AircraftID" validate:"-"`
}
