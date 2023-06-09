package main

import (
	"log"

	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/internal/delivery"
	"github.com/muhammadtaufan/go-sensor-collector/internal/repository"
	"github.com/muhammadtaufan/go-sensor-collector/internal/usecase"
	"github.com/muhammadtaufan/go-sensor-collector/pkg"
)

func main() {
	appConfig := config.LoadConfig()

	db, err := pkg.InitDatabase(appConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	sensorRepo, err := repository.NewSensorRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize sensor repository: %v", err)
	}

	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize sensor repository: %v", err)
	}

	sensorUsecase := usecase.NewSensorUsecase(sensorRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	grpcDelivery := delivery.NewGRPCDelivery(sensorUsecase)
	errChan := make(chan error)

	go func() {
		errChan <- grpcDelivery.RunGRPCServer(appConfig)
	}()

	apiServer := delivery.NewAPIServer(sensorUsecase, userUsecase, appConfig)
	apiServer.StartServer(appConfig)

	err = <-errChan
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
