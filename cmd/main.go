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

	repo, err := repository.NewSensorRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize sensor repository: %v", err)
	}

	usecase := usecase.NewSensorUsecase(repo)
	grpcDelivery := delivery.NewGRPCDelivery(usecase)

	if err := grpcDelivery.RunGRPCServer(appConfig); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
