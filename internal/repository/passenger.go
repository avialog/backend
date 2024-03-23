package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type PassengerRepository interface {
	Save(passenger model.Passenger) (model.Passenger, error)
	GetByID(id uint) (model.Passenger, error)
	GetByFlightID(id uint) ([]model.Passenger, error)
	Update(passenger model.Passenger) (model.Passenger, error)
	DeleteByID(id uint) error
	SaveTx(tx *gorm.DB, passenger model.Passenger) (model.Passenger, error)
	DeleteByFlightIDTx(tx *gorm.DB, flightID uint) error
}

type passenger struct {
	db *gorm.DB
}

func newPassengerRepository(db *gorm.DB) PassengerRepository {
	return &passenger{
		db: db,
	}
}

func (a passenger) Save(passenger model.Passenger) (model.Passenger, error) {
	result := a.db.Create(&passenger)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

// SPECJALNY SAVE POD TRANSAKCJE:
func (a passenger) SaveTx(tx *gorm.DB, passenger model.Passenger) (model.Passenger, error) {
	result := tx.Create(&passenger)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil

}

// /////////////////////////////
func (a passenger) GetByID(id uint) (model.Passenger, error) {
	var passenger model.Passenger
	result := a.db.First(&passenger, id)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

func (a passenger) Update(passenger model.Passenger) (model.Passenger, error) {
	if _, err := a.GetByID(passenger.ID); err != nil {
		return model.Passenger{}, err
	}

	result := a.db.Save(&passenger)
	if result.Error != nil {
		return model.Passenger{}, result.Error
	}

	return passenger, nil
}

func (a passenger) DeleteByID(id uint) error {
	if _, err := a.GetByID(id); err != nil {
		return err
	}

	result := a.db.Delete(&model.Passenger{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a passenger) GetByFlightID(flightID uint) ([]model.Passenger, error) {
	var passengers []model.Passenger
	result := a.db.Where("flight_id = ?", flightID).Find(&passengers)
	if result.Error != nil {
		return nil, result.Error
	}
	return passengers, nil
}

// analogicznie nie sprawdzamy czy istnieje rekord, bo to jest w transakcji
func (a passenger) DeleteByFlightIDTx(tx *gorm.DB, flightID uint) error {
	result := tx.Where("flight_id = ?", flightID).Delete(&model.Passenger{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
