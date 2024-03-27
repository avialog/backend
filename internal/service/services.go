package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

type Services interface {
	Contact() ContactService
	Aircraft() AircraftService
	User() UserService
	Logbook() LogbookService
}

type services struct {
	contactService  ContactService
	aircraftService AircraftService
	userService     UserService
	logbookService  LogbookService
}

func NewServices(repositories repository.Repositories, config dto.Config, validator *validator.Validate) Services {
	contactService := newContactService(repositories.Contact(), config, validator)
	aircraftService := newAircraftService(repositories.Aircraft(), repositories.Flight(), config, validator)
	userService := newUserService(repositories.User(), config)
	logbookService := newLogbookService(repositories.Flight(), repositories.Landing(), repositories.Passenger(), repositories.Aircraft(), config, validator)
	return &services{
		contactService:  contactService,
		aircraftService: aircraftService,
		userService:     userService,
		logbookService:  logbookService,
	}
}

func (s *services) Contact() ContactService {
	return s.contactService
}

func (s *services) Aircraft() AircraftService { return s.aircraftService }

func (s *services) User() UserService { return s.userService }

func (s *services) Logbook() LogbookService { return s.logbookService }
