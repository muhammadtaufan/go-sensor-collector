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
