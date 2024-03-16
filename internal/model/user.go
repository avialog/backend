package model

import "time"

type User struct {
	ID           int64
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	AvatarURL    string
	SignatureURL string
	Country      Country
	Phone        string
	Street       string
	City         string
	Company      string
	Timezone     time.Location
}
