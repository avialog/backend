package dto

type AircraftRequest struct {
	RegistrationNumber string  `json:"registration_number"`
	AircraftModel      string  `json:"aircraft_model"`
	Remarks            *string `json:"remarks"`
	ImageURL           *string `json:"image_url"`
}
