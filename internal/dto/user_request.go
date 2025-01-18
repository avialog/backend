package dto

import "github.com/avialog/backend/internal/model"

type UserRequest struct {
	FirstName     *string        `json:"first_name"`
	LastName      *string        `json:"last_name"`
	AvatarURL     *string        `json:"avatar_url"`
	SignatureURL  *string        `json:"signature_url"`
	Country       *model.Country `json:"country"`
	Phone         *string        `json:"phone"`
	Address       *string        `json:"address"`
	LicenseNumber *string        `json:"license_number"`
	Timezone      *string        `json:"timezone"`
}
