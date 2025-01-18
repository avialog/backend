package dto

import (
	"time"

	"github.com/avialog/backend/internal/model"
)

type LogbookResponse struct {
	FlightID            uint             `json:"flight_id"`
	AircraftID          uint             `json:"aircraft_id"`
	TakeoffTime         time.Time        `json:"takeoff_time"`
	TakeoffAirportCode  string           `json:"takeoff_airport_code"`
	LandingTime         time.Time        `json:"landing_time"`
	LandingAirportCode  string           `json:"landing_airport_code"`
	Style               model.Style      `json:"style"`
	Remarks             *string          `json:"remarks"`
	PersonalRemarks     *string          `json:"personal_remarks"`
	TotalBlockTime      *time.Duration   `json:"total_block_time"`
	PilotInCommandTime  *time.Duration   `json:"pilot_in_command_time"`
	SecondInCommandTime *time.Duration   `json:"second_in_command_time"`
	DualReceivedTime    *time.Duration   `json:"dual_received_time"`
	DualGivenTime       *time.Duration   `json:"dual_given_time"`
	MultiPilotTime      *time.Duration   `json:"multi_pilot_time"`
	NightTime           *time.Duration   `json:"night_time"`
	IFRTime             *time.Duration   `json:"ifr_time"`
	IFRActualTime       *time.Duration   `json:"ifr_actual_time"`
	IFRSimulatedTime    *time.Duration   `json:"ifr_simulated_time"`
	CrossCountryTime    *time.Duration   `json:"cross_country_time"`
	SimulatorTime       *time.Duration   `json:"simulator_time"`
	SignatureURL        *string          `json:"signature_url"`
	FSTDtype            *string          `json:"fstd_type"` // FNPT, FTD, FFS, OTD
	MyRole              model.Role       `json:"my_role"`   // moja rola
	Passengers          []PassengerEntry `json:"passengers"`
	Landings            []LandingEntry   `json:"landings"`
}
