package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
)

type Services interface {
	Contact() ContactService
	Aircraft() AircraftService
}

type services struct {
	contactService  ContactService
	aircraftService AircraftService
}

func NewServices(repositories repository.Repositories, config dto.Config) Services {
	contactService := newContactService(repositories.Contact(), config)
	aircraftService := newAircraftService(repositories.Aircraft(), repositories.Flight(), config)
	return &services{
		contactService:  contactService,
		aircraftService: aircraftService,
	}
}

func (s services) Contact() ContactService {
	return s.contactService
}

func (s services) Aircraft() AircraftService { return s.aircraftService }
