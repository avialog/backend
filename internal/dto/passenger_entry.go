package dto

import "github.com/avialog/backend/internal/model"

type PassengerEntry struct {
	Role         model.Role `json:"role"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Company      string     `json:"company"`
	Phone        string     `json:"phone"`
	EmailAddress string     `json:"email_address"`
	Note         string     `json:"note"`
}
