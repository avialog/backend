package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {

	config := dto.Config{
		DSN: os.Getenv("DSN"),
	}

	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres  password=testserver dbname=newDb sslmode=disable"), &gorm.Config{})

	if err != nil {
		logrus.Panic(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Contact{}, &model.Aircraft{}, &model.Flight{}, &model.Landing{}, &model.Passenger{})
	if err != nil {
		panic(err)
	}

	repositories, err := repository.NewRepositories(db)
	services := service.NewServices(repositories, config)
	_, err = services.Aircraft().InsertAircraft(1, dto.AircraftRequest{
		AircraftModel:      "kolejna",
		RegistrationNumber: "ds",
		ImageURL:           "https://example.com/image.jpg",
		Remarks:            "This is a test aircraft",
	})

	if err != nil {
		logrus.Panic(err)
	}

	//_ = repositories
	//
	//port := "3000"
	//if os.Getenv("PORT") != "" {
	//	port = os.Getenv("PORT")
	//}
	//
	//err = server.Run(":" + port)
	//if err != nil {
	//	logrus.Panic(err)
	//}
}
