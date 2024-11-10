package dto

type AircraftResponse struct {
	ID                 uint    `json:"id"`
	RegistrationNumber string  `json:"registration_number" binding:"required"`
	AircraftModel      string  `json:"aircraft_model" binding:"required"`
	Remarks            *string `json:"remarks"`
	ImageURL           *string `json:"image_url"`
}
