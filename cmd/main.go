package main

import (
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/utils"
	"github.com/go-playground/validator/v10"
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
	_ = service.NewServices(repositories, config, utils.GetValidator())

	x := utils.GetValidator()
	passenger := model.Passenger{
		FlightID:     uint(1),
		Role:         "ds",
		FirstName:    "John",
		LastName:     "Doe",
		Company:      "Company",
		Phone:        "123456789",
		EmailAddress: "esee",
		Note:         "note",
	}
	err = x.Struct(passenger)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("Namespace: ", err.Namespace())
			fmt.Println("Field: ", err.Field()) //TOOOO
			fmt.Println("StructNamespace: ", err.StructNamespace())
			fmt.Println("StructField: ", err.StructField())
			fmt.Println("Tag: ", err.Tag())
			fmt.Println("ActualTag: ", err.ActualTag())
			fmt.Println("Kind: ", err.Kind())
			fmt.Println("Type: ", err.Type())
			fmt.Println("Value: ", err.Value())
			fmt.Println("Param: ", err.Param())
			fmt.Println("pole... required")
		}

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
