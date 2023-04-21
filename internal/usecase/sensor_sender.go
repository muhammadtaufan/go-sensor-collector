package usecase

import (
	"context"
	"log"

	"github.com/muhammadtaufan/go-sensor-collector/internal/repository"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type SensorSender interface {
	SendSensorData(ctx context.Context, data *types.SensorData) error
}

type sensorSender struct {
	repo repository.SensorRepository
}

func NewSensorUsecase(repo repository.SensorRepository) SensorSender {
	return &sensorSender{repo: repo}
}

func (ss *sensorSender) SendSensorData(ctx context.Context, data *types.SensorData) error {
	err := ss.repo.Add(data)
	if err != nil {
		return err
	}

	log.Println("Data is saved")
	return nil
}
