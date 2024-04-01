package controller

import (
	"github.com/avialog/backend/internal/config"
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
	config            config.Config
	contactController ContactController
	authService       service.AuthService
}

func NewControllers(services service.Services, config config.Config) Controllers {
	userController := newUserController(services.User())
	contactController := newContactController(services.Contact())
	infoController := newInfoController()
	authService := services.Auth()
	return &controllers{
		userController:    userController,
		contactController: contactController,
		infoController:    infoController,
		config:            config,
		authService:       authService,
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
	authenticated.Use(middleware.AuthJWT(c.authService))

	authenticated.GET("/profile", c.userController.GetUser)
	authenticated.PUT("/profile", c.userController.UpdateProfile)

	authenticated.GET("/contacts", c.contactController.GetContacts)
	authenticated.POST("/contacts", c.contactController.InsertContact)
	authenticated.PUT("/contacts/:id", c.contactController.UpdateContact)
	authenticated.DELETE("/contacts/:id", c.contactController.DeleteContact)
}
