package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type AircraftRepository interface {
	Save(aircraft model.Aircraft) (model.Aircraft, error)
	GetByID(id uint) (model.Aircraft, error)
	GetByUserID(id uint) ([]model.Aircraft, error)
	Update(aircraft model.Aircraft) (model.Aircraft, error)
	DeleteByID(id uint) error
}

type aircraft struct {
	db *gorm.DB
}

func newAircraftRepository(db *gorm.DB) AircraftRepository {
	return &aircraft{
		db: db,
	}
}

func (a aircraft) Save(aircraft model.Aircraft) (model.Aircraft, error) {
	result := a.db.Create(&aircraft)
	if result.Error != nil {
		return model.Aircraft{}, result.Error
	}

	return aircraft, nil
}

func (a aircraft) GetByID(id uint) (model.Aircraft, error) {
	var aircraft model.Aircraft
	result := a.db.First(&aircraft, id)
	if result.Error != nil {
		return model.Aircraft{}, result.Error
	}
	return aircraft, nil
}

func (a aircraft) Update(aircraft model.Aircraft) (model.Aircraft, error) {
	if _, err := a.GetByID(aircraft.ID); err != nil {
		return model.Aircraft{}, err
	}

	result := a.db.Save(&aircraft)
	if result.Error != nil {
		return model.Aircraft{}, result.Error
	}

	return aircraft, nil
}

func (a aircraft) DeleteByID(id uint) error {
	if _, err := a.GetByID(id); err != nil {
		return err
	}

	result := a.db.Delete(&model.Aircraft{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a aircraft) GetByUserID(userId uint) ([]model.Aircraft, error) {
	var aircraft []model.Aircraft
	result := a.db.Where("user_id = ?", userId).Find(&aircraft)
	if result.Error != nil {
		return nil, result.Error
	}
	return aircraft, nil
}
