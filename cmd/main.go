package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/pkg"
	sensor "github.com/muhammadtaufan/go-sensor/proto"
	"google.golang.org/grpc"
)

type server struct {
	sensor.UnimplementedSensorServiceServer
}

func (s *server) SendSensorData(ctx context.Context, in *sensor.SensorData) (*sensor.SensorDataResponse, error) {
	log.Printf("Received data: %v", in)
	return &sensor.SensorDataResponse{Success: true}, nil
}

func main() {
	appConfig := config.LoadConfig()

	db, err := pkg.InitDatabase(appConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", appConfig.GRPC_HOST, appConfig.GRPC_PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	sensor.RegisterSensorServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
