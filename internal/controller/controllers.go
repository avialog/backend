package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Controllers interface {
	User() UserController
	Info() InfoController
	Route(server *gin.Engine)
}

type controllers struct {
	userController    UserController
	contactController ContactController
	infoController    InfoController
	config            dto.Config
}

func NewControllers(services service.Services, config dto.Config) Controllers {
	userController := newUserController(services.User())
	contactController := newContactController(services.Contact())
	infoController := newInfoController()
	return &controllers{
		userController:    userController,
		contactController: contactController,
		infoController:    infoController,
		config:            config,
	}
}

func (c *controllers) User() UserController {
	return c.userController
}

func (c *controllers) Info() InfoController {
	return c.infoController
}

func (c *controllers) Route(server *gin.Engine) {
	server.GET("/info", c.infoController.Info)
	server.GET("/profile", c.userController.GetProfile)
	server.PUT("/profile", c.userController.UpdateProfile)

	server.GET("/contacts", c.contactController.GetContacts)
	server.POST("/contacts", c.contactController.InsertContact)
	server.PUT("/contacts/:id", c.contactController.UpdateContact)
	server.DELETE("/contacts/:id", c.contactController.DeleteContact)
}
