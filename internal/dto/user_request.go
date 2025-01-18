package dto

import "github.com/avialog/backend/internal/model"

type UserRequest struct {
	FirstName    *string        `json:"first_name"`
	LastName     *string        `json:"last_name"`
	AvatarURL    *string        `json:"avatar_url"`
	SignatureURL *string        `json:"signature_url"`
	Country      *model.Country `json:"country"`
	Phone        *string        `json:"phone"`
	Address      *string        `json:"address"`
	Timezone     *string        `json:"timezone"`
}
