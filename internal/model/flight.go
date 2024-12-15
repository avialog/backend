package model

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	UserID              string      `gorm:"required; not null; default:null" validate:"required"`
	User                User        `validate:"-"`
	AircraftID          uint        `gorm:"required; not null; default:null" validate:"required"`
	Aircraft            Aircraft    `validate:"-"`
	Passengers          []Passenger `gorm:"foreignKey:FlightID" validate:"-"`
	Landings            []Landing   `gorm:"foreignKey:FlightID" validate:"-"`
	TakeoffTime         time.Time   `gorm:"required; not null; default:null" validate:"required"`
	TakeoffAirportCode  string      `gorm:"required; not null; default:null" validate:"required"`
	LandingTime         time.Time   `gorm:"required; not null; default:null" validate:"required"`
	LandingAirportCode  string      `gorm:"required; not null; default:null" validate:"required"`
	Style               Style       `gorm:"required; not null; default:null" validate:"required,style"`
	MyRole              Role        `gorm:"required; not null; default:null" validate:"required,role"`
	Remarks             *string
	PersonalRemarks     *string
	FSTDtype            *string // FNPT, FTD, FFS, OTD
	TotalBlockTime      *time.Duration
	PilotInCommandTime  *time.Duration
	SecondInCommandTime *time.Duration
	DualReceivedTime    *time.Duration
	DualGivenTime       *time.Duration
	MultiPilotTime      *time.Duration
	NightTime           *time.Duration
	IFRTime             *time.Duration
	IFRActualTime       *time.Duration
	IFRSimulatedTime    *time.Duration
	CrossCountryTime    *time.Duration
	SimulatorTime       *time.Duration
	SignatureURL        *string
}
