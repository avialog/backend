package repository

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/infrastructure"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=passenger.go -destination=passenger_mock.go -package repository
type PassengerRepository interface {
	Create(passenger model.Passenger) (model.Passenger, error)
	GetByID(id uint) (model.Passenger, error)
	GetByFlightID(id uint) ([]model.Passenger, error)
	Save(passenger model.Passenger) (model.Passenger, error)
	DeleteByID(id uint) error
	CreateTx(tx infrastructure.Database, passenger model.Passenger) (model.Passenger, error)
	DeleteByFlightIDTx(tx infrastructure.Database, flightID uint) error
}

type passenger struct {
	db *gorm.DB
}

func newPassengerRepository(db *gorm.DB) PassengerRepository {
	return &passenger{
		db: db,
	}
}

func (a *passenger) Create(passenger model.Passenger) (model.Passenger, error) {
	result := a.db.Create(&passenger)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

func (a *passenger) CreateTx(tx infrastructure.Database, passenger model.Passenger) (model.Passenger, error) {
	result := tx.Create(&passenger)
	if result.Error != nil {
		return model.Passenger{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}

	return passenger, nil

}

func (a *passenger) GetByID(id uint) (model.Passenger, error) {
	var passenger model.Passenger
	result := a.db.First(&passenger, id)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

func (a *passenger) Save(passenger model.Passenger) (model.Passenger, error) {
	if _, err := a.GetByID(passenger.ID); err != nil {
		return model.Passenger{}, err
	}

	result := a.db.Save(&passenger)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

func (a *passenger) DeleteByID(id uint) error {
	result := a.db.Delete(&model.Passenger{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("passenger cannot be deleted")
	}

	return nil
}

func (a *passenger) GetByFlightID(flightID uint) ([]model.Passenger, error) {
	var passengers []model.Passenger
	result := a.db.Where("flight_id = ?", flightID).Find(&passengers)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}
	return passengers, nil
}

func (a *passenger) DeleteByFlightIDTx(tx infrastructure.Database, flightID uint) error {
	result := tx.Where("flight_id = ?", flightID).Delete(&model.Passenger{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
