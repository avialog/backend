package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
)

type Services interface {
	Contact() ContactService
	Aircraft() AircraftService
	User() UserService
}

type services struct {
	contactService  ContactService
	aircraftService AircraftService
	userService     UserService
}

func NewServices(repositories repository.Repositories, config dto.Config) Services {
	contactService := newContactService(repositories.Contact(), config)
	aircraftService := newAircraftService(repositories.Aircraft(), repositories.Flight(), config)
	userService := newUserService(repositories.User(), config)
	return &services{
		contactService:  contactService,
		aircraftService: aircraftService,
		userService:     userService,
	}
}

func (s services) Contact() ContactService {
	return s.contactService
}

func (s services) Aircraft() AircraftService { return s.aircraftService }

func (s services) User() UserService { return s.userService }
