package dto

import (
	"time"

	"github.com/avialog/backend/internal/model"
)

type LogbookRequest struct {
	AircraftID          uint             `json:"aircraft_id"`          // used for pdf export
	TakeoffTime         time.Time        `json:"takeoff_time"`         //used for pdf export
	TakeoffAirportCode  string           `json:"takeoff_airport_code"` //used for pdf export
	LandingTime         time.Time        `json:"landing_time"`         //used for pdf export
	LandingAirportCode  string           `json:"landing_airport_code"` //used for pdf export
	Style               model.Style      `json:"style"`
	Remarks             *string          `json:"remarks"` //used for pdf export
	PersonalRemarks     *string          `json:"personal_remarks"`
	TotalBlockTime      *time.Duration   `json:"total_block_time"`       //used for pdf export
	PilotInCommandTime  *time.Duration   `json:"pilot_in_command_time"`  //used for pdf export
	SecondInCommandTime *time.Duration   `json:"second_in_command_time"` //used for pdf export
	DualReceivedTime    *time.Duration   `json:"dual_received_time"`     //used for pdf export
	DualGivenTime       *time.Duration   `json:"dual_given_time"`        //used for pdf export
	MultiPilotTime      *time.Duration   `json:"multi_pilot_time"`
	NightTime           *time.Duration   `json:"night_time"` //used for pdf export
	IFRTime             *time.Duration   `json:"ifr_time"`   //used for pdf export
	IFRActualTime       *time.Duration   `json:"ifr_actual_time"`
	IFRSimulatedTime    *time.Duration   `json:"ifr_simulated_time"` 
	CrossCountryTime    *time.Duration   `json:"cross_country_time"`
	SimulatorTime       *time.Duration   `json:"simulator_time"` // used for pdf export
	SignatureURL        *string          `json:"signature_url"`
	FSTDtype            *string          `json:"fstd_type"` // FNPT, FTD, FFS, OTD // used for pdf export
	MyRole              model.Role       `json:"my_role"`   // moja rola // used for pdf export
	Passengers          []PassengerEntry `json:"passengers"`
	Landings            []LandingEntry   `json:"landings"`
}
