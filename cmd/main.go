package main

import (
	"context"
	"log"
	"net"

	sensor "github.com/muhammadtaufan/go-sensor/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

type server struct {
	sensor.UnimplementedSensorServiceServer
}

func (s *server) SendSensorData(ctx context.Context, in *sensor.SensorData) (*sensor.SensorDataResponse, error) {
	log.Printf("Received data: %v", in)
	return &sensor.SensorDataResponse{Success: true}, nil
}

func main() {
	// Set up and start the gRPC server (Microservice B)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sensor.RegisterSensorServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
