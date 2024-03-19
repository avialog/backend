package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type FlightRepository interface {
	Save(flight model.Flight) (model.Flight, error)
	GetByID(id uint) (model.Flight, error)
	GetByUserID(userID uint) ([]model.Flight, error)
	GetByAircraftID(aircraftID uint) ([]model.Flight, error)
	Update(flight model.Flight) (model.Flight, error)
	DeleteByID(id uint) error
}

type flight struct {
	db *gorm.DB
}

func newFlightRepository(db *gorm.DB) FlightRepository {
	return &flight{
		db: db,
	}
}

func (f flight) Save(flight model.Flight) (model.Flight, error) {
	result := f.db.Create(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f flight) GetByID(id uint) (model.Flight, error) {
	var flight model.Flight
	result := f.db.First(&flight, id)

	if result.Error != nil {
		return model.Flight{}, result.Error
	}
	return flight, nil
}

func (f flight) Update(flight model.Flight) (model.Flight, error) {
	if _, err := f.GetByID(flight.ID); err != nil {
		return model.Flight{}, err
	}

	result := f.db.Save(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f flight) DeleteByID(id uint) error {
	if _, err := f.GetByID(id); err != nil {
		return err
	}

	result := f.db.Delete(&model.Flight{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (f flight) GetByUserID(userID uint) ([]model.Flight, error) {
	var flights []model.Flight

	result := f.db.Where("user_id = ?", userID).Find(&flights)
	if result.Error != nil {
		return nil, result.Error
	}

	return flights, nil
}

func (f flight) GetByAircraftID(aircraftID uint) ([]model.Flight, error) {
	var flights []model.Flight

	result := f.db.Where("aircraft_id = ?", aircraftID).Find(&flights)
	if result.Error != nil {
		return nil, result.Error
	}

	return flights, nil
}
