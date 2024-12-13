package dto

type AircraftRequest struct {
	RegistrationNumber string  `json:"registration_number" binding:"required"` // used for pdf export
	AircraftModel      string  `json:"aircraft_model" binding:"required"`      // used for pdf export
	IsSingleEngine     string  `json:"is_single_engine" binding:"required"`    // used for pdf export
	Remarks            *string `json:"remarks"`
	ImageURL           *string `json:"image_url"`
}
