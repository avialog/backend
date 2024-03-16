package model

type Contact struct {
	ID           int64
	UserID       int64
	AvatarURL    string
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	EmailAddress string
	Note         string
}
