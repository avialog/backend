package repository

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=aircraft.go -destination=aircraft_mock.go -package repository
type AircraftRepository interface {
	Create(aircraft model.Aircraft) (model.Aircraft, error)
	GetByUserIDAndID(userID string, id uint) (model.Aircraft, error)
	GetByUserID(userID string) ([]model.Aircraft, error)
	Save(aircraft model.Aircraft) (model.Aircraft, error)
	DeleteByUserIDAndID(userID string, id uint) error
}

type aircraft struct {
	db *gorm.DB
}

func newAircraftRepository(db *gorm.DB) AircraftRepository {
	return &aircraft{
		db: db,
	}
}

func (a *aircraft) Create(aircraft model.Aircraft) (model.Aircraft, error) {
	result := a.db.Create(&aircraft)
	if result.Error != nil {
		return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}

	return aircraft, nil
}

func (a *aircraft) GetByUserIDAndID(userID string, id uint) (model.Aircraft, error) {
	var aircraft model.Aircraft
	result := a.db.Where("user_id = ? AND id = ?", userID, id).First(&aircraft)
	if result.Error != nil {
		return model.Aircraft{}, result.Error
	}
	return aircraft, nil
}

func (a *aircraft) Save(aircraft model.Aircraft) (model.Aircraft, error) {
	result := a.db.Save(&aircraft)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrNotFound, result.Error)
		}
		return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}

	return aircraft, nil
}

func (a *aircraft) DeleteByUserIDAndID(userID string, id uint) error {
	result := a.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Aircraft{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("aircraft cannot be deleted")
	}

	return nil
}

func (a *aircraft) GetByUserID(userID string) ([]model.Aircraft, error) {
	var aircraft []model.Aircraft
	result := a.db.Where("user_id = ?", userID).Find(&aircraft)
	if result.Error != nil {
		return nil, result.Error
	}
	return aircraft, nil
}
