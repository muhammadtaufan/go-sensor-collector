package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type SensorRepository interface {
	Add(ctx context.Context, data *types.SensorData) error
	GetSensorDataByIDs(ctx context.Context, id1 string, id2 int) ([]types.SensorData, error)
}

type sensorRepository struct {
	db *sqlx.DB
}

func NewSensorRepository(db *sqlx.DB) (SensorRepository, error) {
	return &sensorRepository{
		db: db,
	}, nil
}

func (sr *sensorRepository) Add(ctx context.Context, data *types.SensorData) error {
	query := `INSERT INTO sensor (id, sensor_value, sensor_type, id1, id2, created_at, updated_at)
				VALUES (:id, :sensor_value, :sensor_type, :id1, :id2, :created_at, NOW())`
	_, err := sr.db.NamedExec(query, data)
	if err != nil {
		return err
	}
	return nil
}

func (sr *sensorRepository) GetSensorDataByIDs(ctx context.Context, id1 string, id2 int) ([]types.SensorData, error) {
	query := `SELECT id, sensor_value, sensor_type, id1, id2, created_at from sensor WHERE id1 = ? AND id2 = ?`
	var sensorData []types.SensorData
	err := sr.db.Select(&sensorData, query, id1, id2)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return sensorData, nil
}
