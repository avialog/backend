package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
)

type Services interface {
	Contact() ContactService
}

type services struct {
	contactService ContactService
}

func NewServices(repositories repository.Repositories, config dto.Config) Services {
	contactService := newContactService(repositories.Contact(), config)
	
	return &services{
		contactService: contactService,
	}
}

func (s services) Contact() ContactService {
	return s.contactService
}
