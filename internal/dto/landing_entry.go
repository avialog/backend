package dto

import "github.com/avialog/backend/internal/model"

type LandingEntry struct {
	ApproachType model.ApproachType `json:"approach_type"`
	Count        *uint              `json:"count"`
	NightCount   *uint              `json:"night_count"`
	DayCount     *uint              `json:"day_count"`
	AirportCode  *string            `json:"airport_code"`
}
