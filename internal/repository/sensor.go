package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type SensorRepository interface {
	Add(ctx context.Context, data *types.SensorData) error
	GetSensorData(ctx context.Context, id1 *string, id2 *int, startDate, endDate *time.Time) ([]types.SensorData, error)
	DeleteSensorData(ctx context.Context, id1 *string, id2 *int, startDate, endDate *time.Time) error
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

func (sr *sensorRepository) GetSensorData(ctx context.Context, id1 *string, id2 *int, startDate, endDate *time.Time) ([]types.SensorData, error) {
	query := `SELECT id, sensor_value, sensor_type, id1, id2, created_at from sensor WHERE 1`
	var args []interface{}

	if id1 != nil {
		query += " AND id1 = ?"
		args = append(args, *id1)
	}

	if id2 != nil {
		query += " AND id2 = ?"
		args = append(args, *id2)
	}

	if startDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *endDate)
	}

	var sensorData []types.SensorData
	err := sr.db.Select(&sensorData, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return sensorData, nil
}

func (sr *sensorRepository) DeleteSensorData(ctx context.Context, id1 *string, id2 *int, startDate, endDate *time.Time) error {
	query := `Delete from sensor WHERE 1`
	var args []interface{}

	if id1 != nil {
		query += " AND id1 = ?"
		args = append(args, *id1)
	}

	if id2 != nil {
		query += " AND id2 = ?"
		args = append(args, *id2)
	}

	if startDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *endDate)
	}

	result, err := sr.db.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("no record found")
	}
	return nil
}
