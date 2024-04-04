package model

import "gorm.io/gorm"

type Passenger struct {
	gorm.Model
	FlightID     uint   `gorm:"required; not null; default:null" validate:"required"`
	Flight       Flight `validate:"-"`
	Role         Role   `gorm:"required; not null; default:null" validate:"required,role"`
	FirstName    string `gorm:"required; not null; default:null" validate:"required"`
	LastName     *string
	Company      *string
	Phone        *string
	EmailAddress *string
	Note         *string
}
