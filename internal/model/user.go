package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	FirstName     *string
	LastName      *string
	Email         string `gorm:"required; not null"`
	AvatarURL     *string
	SignatureURL  *string
	Country       *Country
	Phone         *string
	Street        *string
	City          *string
	Company       *string
	Timezone      *string
	LicenseNumber *string
	Contacts      []Contact  `gorm:"foreignKey:UserID" validate:"-"`
	Aircraft      []Aircraft `gorm:"foreignKey:UserID" validate:"-"`
	Flights       []Flight   `gorm:"foreignKey:UserID" validate:"-"`
}
