package types

import "time"

type SensorData struct {
	ID          string    `db:"id"`
	SensorValue float64   `db:"sensor_value"`
	SensorType  string    `db:"sensor_type"`
	ID1         string    `db:"id1"`
	ID2         int       `db:"id2"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type SensorDataResponse struct {
	ID          string    `json:"id"`
	SensorValue float64   `json:"sensor_value"`
	SensorType  string    `json:"sensor_type"`
	ID1         string    `json:"id1"`
	ID2         int       `json:"id2"`
	CreatedAt   time.Time `json:"created_at"`
}

type SensorDataRequest struct {
	ID1       string `json:"id1"`
	ID2       string `json:"id2"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UpdateSensorDataRequest struct {
	SensorValue float64 `json:"sensor_value"`
}
