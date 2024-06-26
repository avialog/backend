package service

import (
	authV4 "firebase.google.com/go/v4/auth"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

type Services interface {
	Contact() ContactService
	Aircraft() AircraftService
	User() UserService
	Logbook() LogbookService
	Auth() AuthService
}

type services struct {
	contactService  ContactService
	aircraftService AircraftService
	userService     UserService
	logbookService  LogbookService
	authService     AuthService
}

func NewServices(repositories repository.Repositories, config config.Config, validator *validator.Validate, authClient *authV4.Client) Services {
	contactService := newContactService(repositories.Contact(), config, validator)
	aircraftService := newAircraftService(repositories.Aircraft(), repositories.Flight(), config, validator)
	userService := newUserService(repositories.User(), config)
	logbookService := newLogbookService(repositories.Flight(), repositories.Landing(), repositories.Passenger(), repositories.Aircraft(), config, validator)
	authService := newAuthService(repositories.User(), authClient, authV4.IsIDTokenExpired)
	return &services{
		contactService:  contactService,
		aircraftService: aircraftService,
		userService:     userService,
		logbookService:  logbookService,
		authService:     authService,
	}
}

func (s *services) Contact() ContactService {
	return s.contactService
}

func (s *services) Aircraft() AircraftService { return s.aircraftService }

func (s *services) User() UserService { return s.userService }

func (s *services) Logbook() LogbookService { return s.logbookService }

func (s *services) Auth() AuthService { return s.authService }
