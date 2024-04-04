package dto

type ContactRequest struct {
	AvatarURL    *string `json:"avatar_url"`
	FirstName    string  `json:"first_name" binding:"required"`
	LastName     *string `json:"last_name"`
	Company      *string `json:"company"`
	Phone        *string `json:"phone"`
	EmailAddress *string `json:"email_address"`
	Note         *string `json:"note"`
}
