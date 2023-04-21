package repository

import (
	"database/sql"

	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type SensorRepository interface {
	Add(data *types.SensorData) error
}

type sensorRepository struct {
	db *sql.DB
}

func NewSensorRepository(db *sql.DB) (SensorRepository, error) {
	return &sensorRepository{
		db: db,
	}, nil
}

func (sr *sensorRepository) Add(data *types.SensorData) error {
	query := `INSERT INTO sensor (id, sensor_value, sensor_type, id1, id2, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err := sr.db.Exec(query, data.ID, data.SensorValue, data.SensorType, data.ID1, data.ID2, data.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
