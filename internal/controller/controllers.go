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
	authMiddleware    gin.HandlerFunc
}

func NewControllers(services service.Services, config config.Config) Controllers {
	userController := newUserController(services.User())
	contactController := newContactController(services.Contact())
	infoController := newInfoController()
	authMiddleware := middleware.AuthJWT(services.Auth())
	return &controllers{
		userController:    userController,
		contactController: contactController,
		infoController:    infoController,
		config:            config,
		authMiddleware:    authMiddleware,
	}
}

func (c *controllers) User() UserController {
	return c.userController
}

func (c *controllers) Info() InfoController {
	return c.infoController
}

func (c *controllers) Route(server *gin.Engine) {

	api := server.Group("/api")
	{
		info := api.Group("/info")
		{
			info.GET("", c.infoController.Info)
		}

		authenticated := api.Group("/")
		{
			authenticated.Use(c.authMiddleware)

			profile := authenticated.Group("/profile")
			{
				profile.GET("", c.userController.GetUser)
				profile.PUT("", c.userController.UpdateProfile)
			}

			contacts := authenticated.Group("/contacts")
			{
				contacts.GET("", c.contactController.GetContacts)
				contacts.POST("", c.contactController.InsertContact)
				contacts.PUT(":id", c.contactController.UpdateContact)
				contacts.DELETE(":id", c.contactController.DeleteContact)
			}
		}

	}
}
