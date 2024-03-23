package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var logbookRequest dto.LogbookRequest = dto.LogbookRequest{
	TakeoffTime:        time.Now(),
	TakeoffAirportCode: "PoprawneTransakcjeV2",
	LandingTime:        time.Now().Add(1 * time.Hour),
	LandingAirportCode: "PoprawneTransakcjeV2",
	Style:              "PoprawneTransakcjeV2",
	Remarks:            "PoprawneTransakcjeV2",
	PersonalRemarks:    "PoprawneTransakcjeV2",
	TotalBlockTime:     3 * time.Hour,
	PilotInCommandTime: 2 * time.Hour,
	Passengers: []dto.PassengerEntry{
		{
			Role:         "PoprawneTransakcjeV2",
			FirstName:    "John",
			LastName:     "Doe",
			Company:      "ABC Company",
			Phone:        "123-456-7890",
			EmailAddress: "john.doe@example.com",
			Note:         "Some note",
		},
		{
			Role:         "PoprawneTransakcjeV2",
			FirstName:    "XXX",
			LastName:     "XXX",
			Company:      "XXXX",
			Phone:        "XXX",
			EmailAddress: "XXXX",
			Note:         "XXX",
		},
		{
			Role:         "PoprawneTransakcjeV2",
			FirstName:    "XXX",
			LastName:     "XXX",
			Company:      "XXXX",
			Phone:        "XXX",
			EmailAddress: "XXXX",
			Note:         "XXX",
		},
	},
	Landings: []dto.LandingEntry{
		{
			ApproachType: "PoprawneTransakcjeV2",
			Count:        1,
			NightCount:   0,
			DayCount:     1,
			AirportCode:  "YY",
		},
		{
			ApproachType: "PoprawneTransakcjeV2",
			Count:        1,
			NightCount:   0,
			DayCount:     1,
			AirportCode:  "YYY",
		},
	},
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	logrus.Panic("Error loading .env file")
	//}
	//
	config := dto.Config{
		DSN: os.Getenv("DSN"),
	}

	//db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=testserver dbname=newDb port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}

	//server := gin.Default()

	repositories, err := repository.NewRepositories(db)
	if err != nil {
		logrus.Panic(err)
	}
	services := service.NewServices(repositories, config)

	_ = services
	_ = repositories
	lservice := services.Logbook()
	//err = lservice.InsertLogbookEntry(1, 2, logbookRequest)
	err = lservice.DeleteLogbookEntry(1, 12)
	if err != nil {
		logrus.Panic(err)
	}
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
