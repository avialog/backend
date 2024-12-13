package controller

import (
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/middleware"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Controllers interface {
	User() UserController
	Contact() ContactController
	Logbook() LogbookController
	Info() InfoController
	Route(server *gin.Engine)
	Aircraft() AircraftController
}

type controllers struct {
	userController     UserController
	infoController     InfoController
	config             config.Config
	contactController  ContactController
	authMiddleware     gin.HandlerFunc
	aircraftController AircraftController
	logbookController  LogbookController
}

func NewControllers(services service.Services, config config.Config) Controllers {
	userController := newUserController(services.User())
	contactController := newContactController(services.Contact())
	aircraftController := newAircraftController(services.Aircraft())
	infoController := newInfoController()
	authMiddleware := middleware.AuthJWT(services.Auth())
	flightController := newLogbookController(services.Logbook())
	return &controllers{
		userController:     userController,
		contactController:  contactController,
		infoController:     infoController,
		config:             config,
		authMiddleware:     authMiddleware,
		aircraftController: aircraftController,
		logbookController:  flightController,
	}
}

func (c *controllers) User() UserController {
	return c.userController
}

func (c *controllers) Info() InfoController {
	return c.infoController
}

func (c *controllers) Aircraft() AircraftController { return c.aircraftController }

func (c *controllers) Contact() ContactController { return c.contactController }

func (c *controllers) Logbook() LogbookController { return c.logbookController }

func (c *controllers) Route(server *gin.Engine) {

	server.GET("/healthz", c.infoController.Info)

	api := server.Group("/api")
	{

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

			flights := authenticated.Group("/logbook")
			{
				flights.GET("", c.logbookController.GetLogbookEntries)
				flights.POST("", c.logbookController.InsertLogbookEntry)
				flights.PUT(":id", c.logbookController.UpdateLogbookEntry)
				flights.DELETE(":id", c.logbookController.DeleteLogbookEntry)
				flights.GET("/download", c.logbookController.DownloadLogbookPDF)
			}
			aircraft := authenticated.Group("/aircraft")
			{
				aircraft.GET("", c.aircraftController.GetAircraft)
				aircraft.POST("", c.aircraftController.InsertAircraft)
				aircraft.PUT(":id", c.aircraftController.UpdateAircraft)
				aircraft.DELETE(":id", c.aircraftController.DeleteAircraft)
			}

		}

	}
}
