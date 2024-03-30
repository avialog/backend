package controller

import (
	"firebase.google.com/go/auth"
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
	userController UserController
	infoController InfoController
	config         dto.Config
	authClient     *auth.Client
}

func NewControllers(services service.Services, config dto.Config, authClient *auth.Client) Controllers {
	userController := newUserController(services.User())
	infoController := newInfoController()
	return &controllers{
		userController: userController,
		infoController: infoController,
		config:         config,
		authClient:     authClient,
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
	authenticated.Use(middleware.AuthJWT(c.authClient))

	authenticated.GET("/profile", c.userController.GetUser)
	authenticated.PUT("/profile", c.userController.UpdateProfile)

}
