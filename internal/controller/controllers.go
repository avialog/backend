package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/middleware"
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
	infoController    InfoController
	config            dto.Config
	contactController ContactController
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

	authenticated := server.Group("/")
	authenticated.Use(middleware.AuthJWT(middleware.AuthJWT(c.AuthService)))

	authenticated.GET("/profile", c.userController.GetUser)
	authenticated.PUT("/profile", c.userController.UpdateProfile)

	server.GET("/contacts", c.contactController.GetContacts)
	server.POST("/contacts", c.contactController.InsertContact)
	server.PUT("/contacts/:id", c.contactController.UpdateContact)
	server.DELETE("/contacts/:id", c.contactController.DeleteContact)
}
