package delivery

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
	"github.com/muhammadtaufan/go-sensor-collector/internal/usecase"
	sensorProto "github.com/muhammadtaufan/go-sensor/proto"
	"google.golang.org/grpc"
)

type grpcDelivery struct {
	sensorProto.UnimplementedSensorServiceServer
	usecase usecase.SensorSender
}

func NewGRPCDelivery(usecase usecase.SensorSender) *grpcDelivery {
	return &grpcDelivery{usecase: usecase}
}

func (gd *grpcDelivery) RunGRPCServer(cfg *config.Config) error {
	address := fmt.Sprintf("%s:%s", cfg.GRPC_HOST, cfg.GRPC_PORT)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	sensorProto.RegisterSensorServiceServer(grpcServer, gd)

	log.Printf("gRPC server is running on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("grpc failed to serve: %v", err)
	}

	return nil
}

func (gd *grpcDelivery) SendSensorData(ctx context.Context, in *sensorProto.SensorData) (*sensorProto.SensorDataResponse, error) {
	log.Printf("Received data: %v", in)

	err := gd.usecase.AddSensorData(ctx, &types.SensorData{
		ID:          uuid.New().String(),
		SensorValue: float64(in.SensorValue),
		SensorType:  in.SensorType,
		ID1:         in.ID1,
		ID2:         int(in.ID2),
		CreatedAt:   time.Unix(in.Timestamp, 0).UTC(),
	})

	if err != nil {
		return &sensorProto.SensorDataResponse{Success: false}, err
	}

	return &sensorProto.SensorDataResponse{Success: true}, nil
}
