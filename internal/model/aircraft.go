package model

type Aircraft struct {
	ID                 int64
	UserID             int64
	RegistrationNumber string
	Model              string
	Remarks            string
	ImageURL           string
}
