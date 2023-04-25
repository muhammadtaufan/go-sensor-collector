package usecase

import (
	"context"
	"log"
	"time"

	"github.com/muhammadtaufan/go-sensor-collector/internal/repository"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type SensorSender interface {
	AddSensorData(ctx context.Context, data *types.SensorData) error
	GetSensorData(ctx context.Context, id1 *string, id2 *int, startTime, endTime *time.Time) ([]types.SensorDataResponse, error)
}

type sensorSender struct {
	repo repository.SensorRepository
}

func NewSensorUsecase(repo repository.SensorRepository) SensorSender {
	return &sensorSender{repo: repo}
}

func (ss *sensorSender) AddSensorData(ctx context.Context, data *types.SensorData) error {
	err := ss.repo.Add(ctx, data)
	if err != nil {
		return err
	}

	log.Println("Data is saved")
	return nil
}

func (ss *sensorSender) GetSensorData(ctx context.Context, id1 *string, id2 *int, startTime, endTime *time.Time) ([]types.SensorDataResponse, error) {
	results, err := ss.repo.GetSensorData(ctx, id1, id2, startTime, endTime)
	if err != nil {
		return nil, err
	}

	var response []types.SensorDataResponse
	for _, result := range results {
		response = append(response, types.SensorDataResponse{
			ID:          result.ID,
			SensorValue: result.SensorValue,
			SensorType:  result.SensorType,
			ID1:         result.ID1,
			ID2:         result.ID2,
			CreatedAt:   result.CreatedAt,
		})
	}
	return response, nil
}
