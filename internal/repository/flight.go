package repository

import (
	"errors"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -source=flight.go -destination=flight_mock.go -package repository
type FlightRepository interface {
	Save(flight model.Flight) (model.Flight, error)
	GetByID(id uint) (model.Flight, error)
	GetByUserID(userID uint) ([]model.Flight, error)
	GetByAircraftID(aircraftID uint) ([]model.Flight, error)
	Update(flight model.Flight) (model.Flight, error)
	DeleteByID(id uint) error
	CountByAircraftID(userID, aircraftID uint) (int64, error)
	GetByUserIDAndDate(userID uint, start, end time.Time) ([]model.Flight, error)
	Begin() *gorm.DB
	SaveTx(tx *gorm.DB, flight model.Flight) (model.Flight, error)
	DeleteByIDTx(tx *gorm.DB, id uint) error
	GetByIDTx(tx *gorm.DB, id uint) (model.Flight, error)
	UpdateTx(tx *gorm.DB, flight model.Flight) (model.Flight, error)
}

type flight struct {
	db *gorm.DB
}

func newFlightRepository(db *gorm.DB) FlightRepository {
	return &flight{
		db: db,
	}
}

func (f *flight) Begin() *gorm.DB {
	return f.db.Begin()
}

func (f *flight) Save(flight model.Flight) (model.Flight, error) {
	result := f.db.Create(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f *flight) SaveTx(tx *gorm.DB, flight model.Flight) (model.Flight, error) {
	result := tx.Create(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f *flight) GetByID(id uint) (model.Flight, error) {
	var flight model.Flight
	result := f.db.First(&flight, id)

	if result.Error != nil {
		return model.Flight{}, result.Error
	}
	return flight, nil
}

func (f *flight) Update(flight model.Flight) (model.Flight, error) {
	if _, err := f.GetByID(flight.ID); err != nil {
		return model.Flight{}, err
	}

	result := f.db.Save(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f *flight) UpdateTx(tx *gorm.DB, flight model.Flight) (model.Flight, error) {
	if _, err := f.GetByIDTx(tx, flight.ID); err != nil { //is this necessary?
		return model.Flight{}, err
	}

	result := tx.Save(&flight)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}

func (f *flight) DeleteByID(id uint) error {
	result := f.db.Delete(&model.Flight{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("flight cannot be deleted")
	}
	return nil
}

func (f *flight) GetByUserID(userID uint) ([]model.Flight, error) {
	var flights []model.Flight

	result := f.db.Where("user_id = ?", userID).Find(&flights)
	if result.Error != nil {
		return nil, result.Error
	}

	return flights, nil
}

func (f *flight) GetByAircraftID(aircraftID uint) ([]model.Flight, error) {
	var flights []model.Flight

	result := f.db.Where("aircraft_id = ?", aircraftID).Find(&flights)
	if result.Error != nil {
		return nil, result.Error
	}

	return flights, nil
}

func (f *flight) CountByAircraftID(userID, aircraftID uint) (int64, error) {
	var count int64
	result := f.db.Model(&model.Flight{}).Where("aircraft_id = ? AND user_id = ?", aircraftID, userID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (f *flight) GetByUserIDAndDate(userID uint, start, end time.Time) ([]model.Flight, error) {
	var flights []model.Flight

	result := f.db.Where("user_id = ? AND takeoff_time >= ? AND takeoff_time <= ?", userID, start, end).Find(&flights)
	if result.Error != nil {
		return nil, result.Error
	}

	return flights, nil
}

func (f *flight) DeleteByIDTx(tx *gorm.DB, id uint) error {
	result := tx.Delete(&model.Flight{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("flight cannot be deleted")
	}
	return nil
}

func (f *flight) GetByIDTx(tx *gorm.DB, id uint) (model.Flight, error) {
	var flight model.Flight
	result := tx.First(&flight, id)
	if result.Error != nil {
		return model.Flight{}, result.Error
	}

	return flight, nil
}
