package repository

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/infrastructure"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=landing.go -destination=landing_mock.go -package repository
type LandingRepository interface {
	Create(landing model.Landing) (model.Landing, error)
	GetByID(id uint) (model.Landing, error)
	GetByFlightID(flightID uint) ([]model.Landing, error)
	Save(landing model.Landing) (model.Landing, error)
	DeleteByID(id uint) error
	CreateTx(tx infrastructure.Database, landing model.Landing) (model.Landing, error)
	DeleteByFlightIDTx(tx infrastructure.Database, flightID uint) error
}

type landing struct {
	db *gorm.DB
}

func newLandingRepository(db *gorm.DB) LandingRepository {
	return &landing{
		db: db,
	}
}

func (l *landing) Create(landing model.Landing) (model.Landing, error) {
	result := l.db.Create(&landing)
	if result.Error != nil {
		return model.Landing{}, result.Error
	}

	return landing, nil
}

func (l *landing) CreateTx(tx infrastructure.Database, landing model.Landing) (model.Landing, error) {
	result := tx.Create(&landing)
	if result.Error != nil {
		return model.Landing{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}

	return landing, nil
}

func (l *landing) GetByID(id uint) (model.Landing, error) {
	var landing model.Landing
	result := l.db.First(&landing, id)
	if result.Error != nil {
		return model.Landing{}, result.Error
	}
	return landing, nil
}

func (l *landing) Save(landing model.Landing) (model.Landing, error) {
	if _, err := l.GetByID(landing.ID); err != nil {
		return model.Landing{}, err
	}

	result := l.db.Save(&landing)
	if result.Error != nil {
		return model.Landing{}, result.Error
	}

	return landing, nil
}

func (l *landing) DeleteByID(id uint) error {
	result := l.db.Delete(&model.Landing{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("landing cannot be deleted")
	}

	return nil
}

func (l *landing) GetByFlightID(flightID uint) ([]model.Landing, error) {
	var landings []model.Landing
	result := l.db.Where("flight_id = ?", flightID).Find(&landings)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}
	return landings, nil
}

func (l *landing) DeleteByFlightIDTx(tx infrastructure.Database, flightID uint) error {
	result := tx.Delete(&model.Landing{}, "flight_id = ?", flightID)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", dto.ErrInternalFailure, result.Error)
	}

	return nil
}
